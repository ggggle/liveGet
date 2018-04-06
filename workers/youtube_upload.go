package workers

import (
    "os/exec"
    "os"
    "github.com/ggggle/liveGet/api"
    "fmt"
    "bytes"
    "strings"
    "time"
)

var MIN_FILE_SIZE int64 = 1024 * 128

func touchVideo(fPath string) (int) {
    cmd := exec.Command("ffmpeg", "-i", fPath)
    w := bytes.NewBuffer(nil)
    cmd.Stderr = w
    cmd.Run()
    //不存在该子串说明不是有效的文件
    stream0Exist := strings.Index(string(w.Bytes()), "Stream #0:0")
    if -1 == stream0Exist {
        return -1
    }
    //不存在该子串说明是纯音频
    stream1Exist := strings.Index(string(w.Bytes()), "Stream #0:1")
    if -1 == stream1Exist {
        return 0
    }
    return 1 //正常的视频
}

func combineAudio(fPath string) (output string, ret int) {
    output = fPath + ".mp4"
    cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "/usr/sy_cr.jpg", "-i", fPath,
        "-c:a", "copy", "-c:v", "libx264", "-shortest", output)
    w := bytes.NewBuffer(nil)
    cmd.Stderr = w
    cmd.Run()
    _, err := os.Stat(output)
    if nil!=err{
        ret = -1  //combine失败，记录日志
        api.Logger.Println(string(w.Bytes()))
    } else {
        ret = 0
        os.Remove(fPath)
    }
    return
}

//fPath文件路径  retry失败重试次数
func YoutubeUpload(API *api.LuzhiboAPI, fPath string, retry int) {
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
    if videoRet := touchVideo(fPath); -1 == videoRet {
        return
    } else if 0 == videoRet { //纯音频，无法直接上传到youtube，需要加一张图片合成视频
        combineFile, ret := combineAudio(fPath)
        if -1 == ret{
            return
        }
        fPath = combineFile
    }
    extraInfo, _ := API.G.GetExtraInfo(API.Id)
    site := fmt.Sprintf("[%s]", API.G.Site())
    roomId := API.Id
    //平台-主播名-直播标题-房间id-'结束日期-结束时间'
    title := fmt.Sprintf("%s-%s-%s-%s-%s", site, extraInfo.OwnerName,
        extraInfo.RoomTitle, roomId, info.ModTime().Format("20060102-1504"))
    for ; retry >= 0; retry-- {
        cmd := exec.Command("youtube-upload", "--client-secrets", "/root/.client_secret.json",
            "--privacy", "private", "--title", title, "--playlist", roomId, fPath)
        w := bytes.NewBuffer(nil)
        cmd.Stderr = w
        cmd.Run()
        uploadRet := string(w.Bytes())
        success := strings.Contains(uploadRet, "Video URL")
        if success {
            api.Logger.Printf("[%s]上传成功", fPath)
            os.Remove(fPath)
            return
        } else {
            api.Logger.Print(uploadRet)
            if retry <= 0 {
                tmp := fmt.Sprintf("youtube-upload --client-secrets /root/.client_secret.json --privacy private --title %s --playlist %s %s",
                    title, roomId, fPath)
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
