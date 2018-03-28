package workers

import (
	"github.com/ggggle/luzhibo/api"
	"errors"
	"github.com/ggggle/luzhibo/api/getters"
)

//普通模式

type singleworker struct {
	url      string
	filePath string
	cb       WorkCompletedCallBack
	run      bool
	ch       chan bool
	dl       *downloader
	API      *api.LuzhiboAPI
}

//NewSingleWorker 创建对象
func NewSingleWorker(oa *api.LuzhiboAPI, filepath string, callbcak WorkCompletedCallBack) (r *singleworker, err error) {
	if oa != nil {
		r = &singleworker{}
		var live bool
		_, live, err = oa.GetRoomInfo()
		if err != nil {
			err = errors.New("没有这个房间")
			return
		}
		if !live {
			err = errors.New("房间未开播")
			return
		}
		var t getters.LiveInfo
		t, err = oa.GetLiveInfo()
		if err != nil {
			err = errors.New("获取直播信息失败")
			return
		}
		r.API = oa
		r.url = t.VideoURL
		r.cb = callbcak
		r.filePath = filepath
		return
	}
	err = errors.New("-1") //参数错误
	return
}

//Start 实现接口
func (i *singleworker) Start() {
	if i.run {
		return
	}
	i.run = true
	i.dl = newDownloader(i.url, i.filePath, i.dwnloaderCallback)
	i.ch = make(chan bool, 0)
	if i.API.NeedFFmpeg {
		i.dl.UseFFmpeg()
	}
	i.dl.Start()
}

//Stop 实现接口
func (i *singleworker) Stop() {
	if i.run {
		i.run = false
		if _, r, _, _, _ := i.dl.GetTaskInfo(false); r {
			i.dl.Stop()
		}
		<-i.ch
		close(i.ch)
	}
}

//Restart 实现接口
func (i *singleworker) Restart() (Worker, error) {
	if i.run {
		i.Stop()
	}
	r, e := NewSingleWorker(i.API, i.filePath, i.cb)
	if e == nil {
		i = r
		i.Start()
	}
	return i, e
}

//GetTaskInfo 实现接口
func (i *singleworker) GetTaskInfo(g bool) (int64, bool, int64, string, *getters.LiveInfo) {
	if g {
		r, e := i.API.GetLiveInfo()
		if e == nil {
			return 1, i.run, 1, i.filePath, &r
		}
	}
	return 1, i.run, 1, i.filePath, nil
}

func (i *singleworker) dwnloaderCallback(x int64) {
	if !i.run {
		i.ch <- true
	}
	i.run = false
	if i.cb != nil {
		i.cb(x)
	}
}
