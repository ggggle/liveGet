package workers

import (
    "os/exec"
    "fmt"
    "bytes"
    "strings"
    "github.com/ggggle/liveGet/api"
    "os"
    "time"
)

var ROOT_DIR = "liveGet"

func GetDirID(dirName string) (ID string, exist bool) {
    queryArg := fmt.Sprintf("gdrive list --query \"name='%s'\"", dirName)
    cmd := exec.Command("/bin/sh", "-c", queryArg)
    w := bytes.NewBuffer(nil)
    cmd.Stdout = w
    cmd.Run()
    queryRet := string(w.Bytes())
    /*
    Id                                  Name        Type   Size   Created
    1Y1kGrgQzsDslISYBBgS2Rfg7QTwE_aEe   [斗鱼直播]196   dir           2018-06-07 21:17:27
*/
    lines := strings.Split(queryRet, "\n")
    api.Logger.Println(lines)
    if len(lines) == 1 {
        return
    }
    for _, oneLine := range lines {
        if strings.Contains(oneLine, " dir ") {
            exist = true
            lineSplit := strings.Split(oneLine, " ")
            ID = lineSplit[0]
            return
        }
    }
    return
}

func MakeDir(dirName string, parent ...string) (id string) {
    var args []string
    args = append(args, "mkdir")
    if len(parent)!=0{
        args = append(args, "-p")
        args = append(args, parent[0])
        args = append(args, dirName)
    } else {
        args = append(args, dirName)
    }
    cmd := exec.Command("gdrive", args...)
    w := bytes.NewBuffer(nil)
    cmd.Stdout = w
    cmd.Run()

    ret := string(w.Bytes())
    return strings.Split(ret, " ")[1]
}

//若已经存在该dir，则返回id，否则创建该dir
func DirID(dirName string) (id string)  {
    id, exist := GetDirID(dirName)
    if exist{
        return
    }
    liveGetRootid, exist := GetDirID(ROOT_DIR)
    if !exist{
        liveGetRootid = MakeDir(ROOT_DIR)
    }
    return MakeDir(dirName, liveGetRootid)
}

func GdriveUpload(API *api.LuzhiboAPI, fPath string, retry int) {
    info, err := os.Stat(fPath)
    if err != nil {
        api.Logger.Print(fPath + " error")
        return
    }
    //128KB
    if info.Size() < MIN_FILE_SIZE {
        api.Logger.Printf("[%s]长度太短[%d]", fPath, info.Size())
        return
    }
    if videoRet := TouchVideo(fPath); -1 == videoRet {
        return
    } else if 0 == videoRet { //纯音频，无法直接上传到youtube，需要加一张图片合成视频
        combineFile, ret := CombineAudio(fPath)
        if -1 == ret{
            return
        }
        fPath = combineFile
    }
    extraInfo, _ := API.G.GetExtraInfo(API.Id)
    site := fmt.Sprintf("[%s]", API.G.Site())
    roomId := API.Id
    //平台-主播名-直播标题-房间id-'结束日期-结束时间'
    fileName := fmt.Sprintf("%s-%s-%s-%s-%s.flv", site, extraInfo.OwnerName,
        extraInfo.RoomTitle, roomId, info.ModTime().Format("20060102-1504"))
    parentID := DirID(site + roomId)
    for ; retry >= 0; retry-- {
        cmd := exec.Command("gdrive", "upload", "-p", parentID, "--name", fileName,"--delete", fPath)
        w := bytes.NewBuffer(nil)
        cmd.Stdout = w
        cmd.Run()
        uploadRet := string(w.Bytes())
        success := strings.Contains(uploadRet, "Removed")
        if success {
            api.Logger.Printf("[%s]上传成功", fPath)
            return
        } else {
            api.Logger.Print(uploadRet)
            if retry <= 0 {
                tmp := fmt.Sprintf("gdrive -upload -p %s --name %s --delete %s", parentID, fileName, fPath)
                api.Logger.Printf("上传cmd{%s}", tmp)
                return
            }
            api.Logger.Printf("5min后重试一次，剩余重传次数[%d]", retry)
            select {
            case <-time.After(5 * time.Minute):
            }
        }
    }
}