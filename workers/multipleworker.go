package workers

import (
    "errors"
    "fmt"
    "os"
    "time"
    "github.com/ggggle/luzhibo/api"
    "github.com/ggggle/luzhibo/api/getters"
    "strings"
)

//循环模式

type multipleworker struct {
    dirPath string
    index   int64
    cb      WorkCompletedCallBack
    run     bool
    ch      chan bool
    ch2     chan bool //一次循环结束
    ch3     chan bool //停止
    API     *api.LuzhiboAPI
    sw      *singleworker
}

//NewMultipleWorker 创建对象
func NewMultipleWorker(oa *api.LuzhiboAPI, dirpath string, callbcak WorkCompletedCallBack) (r *multipleworker, err error) {
    if oa != nil {
        r = &multipleworker{}
        _, _, err = oa.GetRoomInfo()
        if err != nil {
            err = errors.New("没有这个房间")
            return
        }
        r.cb = callbcak
        r.dirPath = dirpath
        r.API = oa
        return
    }
    err = errors.New("-1") //参数错误
    return
}

//Start 实现接口
func (i *multipleworker) Start() {
    if i.run {
        return
    }
    i.run = true
    i.ch = make(chan bool, 0)
    i.ch3 = make(chan bool, 1)
    go i.do()
}

//Stop 实现接口
func (i *multipleworker) Stop() {
    if i.run {
        i.run = false
        if i.sw != nil {
            if _, r, _, _, _ := i.sw.GetTaskInfo(false); r {
                i.sw.Stop()
            }
        }
        i.ch3 <- true
        <-i.ch
        close(i.ch)
        close(i.ch3)
    }
}

//Restart 实现接口
func (i *multipleworker) Restart() (Worker, error) {
    if i.run {
        i.Stop()
    }
    r, e := NewMultipleWorker(i.API, i.dirPath, i.cb)
    if e == nil {
        i = r
        i.Start()
    }
    return i, e
}

//GetTaskInfo 实现接口
func (i *multipleworker) GetTaskInfo(g bool) (int64, bool, int64, string, *getters.LiveInfo) {
    if 0 == strings.Compare(i.API.Site, "斗鱼直播") {
        var r getters.LiveInfo
        extraInfo, err := i.API.G.GetExtraInfo(i.API.Id)
        if nil == err {
            r.RoomTitle = extraInfo.RoomTitle
            return 2, i.run, i.index, i.dirPath, &r
        }
    }
    if i.sw != nil {
        _, _, _, _, r := i.sw.GetTaskInfo(g)
        return 2, i.run, i.index, i.dirPath, r
    }
    return 2, i.run, i.index, i.dirPath, nil
}

func (i *multipleworker) do() {
    var ec int64
    var fn string
    for i.run {
        ec = 0
        i.ch2 = make(chan bool, 0)
        i.index++
        fn = fmt.Sprintf("%s/%d.%s", i.dirPath, i.index, i.API.FileExt)
        r, err := NewSingleWorker(i.API, fn, func(x int64) {
            ec = x
            i.ch2 <- true
        })
        b := false
        if err == nil {
            i.sw = r
            b = true
        } else {
            i.index--
            api.Logger.Print(err.Error() + " " + i.API.Id)
        }
        if b {
            i.sw.Start()
            <-i.ch2
            p, err := os.Stat(fn)
            if err == nil {
                if !p.IsDir() && p.Size() == 0 {
                    i.index--
                }
            } else {
                i.index--
            }
        }
        api.Logger.Printf("[%s]err code[%d] index[%d]", i.API.Id, ec, i.index)
        if ec == 5 {
            api.Logger.Print("ec code5 " + i.API.Id)
            break
        }
        if i.run {
            if 4 == ec || 6 == ec || 7 == ec {
                go YoutubeUpload(i.API, fn, 3)
                switch ec {
                case 4: //get下载过程中网络问题导致断开
                    api.Logger.Print("下载数据错误,立即重试")
                    continue
                case 6: //分段下载保存为多个文件
                    api.Logger.Print("next")
                    continue
                case 7:
                    api.Logger.Print("[EVENT]跳过等待")
                    continue
                }
            }
            select {
            case <-i.ch3:
            case <-time.After(1 * time.Minute):
            }
        }
    }
    if !i.run {
        go YoutubeUpload(i.API, fn, 3)
        i.ch <- true
    }
    i.run = false
    if i.cb != nil {
        api.Logger.Print("执行回调")
        i.cb(ec)
    }
}
