package getters

import (
	"errors"
	"regexp"
	"strings"
	"fmt"
)

//panda 熊猫星颜
type xingyan struct{}

//Site 实现接口
func (i *xingyan) Site() string { return "熊猫星颜" }

func (i *xingyan) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *xingyan) SiteURL() string {
	return "http://xingyan.panda.tv"
}

//SiteIcon 实现接口
func (i *xingyan) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *xingyan) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *xingyan) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *xingyan) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("xingyan\\.panda\\.tv/(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://m.api.xingyan.panda.tv/room/baseinfo?xid=" + id
	tmp, err := httpGet(url)
	if strings.Contains(tmp, "\"errno\":0") {
		live = strings.Contains(tmp, "\"playstatus\":\"1\"")
	} else {
		id = ""
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *xingyan) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://m.api.xingyan.panda.tv/room/baseinfo?xid=" + id
	tmp, err := httpGet(url)
	json := pruseJSON(tmp).JToken("data")
	roomInfo, videoInfo, hostInfo := *(json.JToken("roominfo")), *(json.JToken("videoinfo")), *(json.JToken("hostinfo"))
	nick := fmt.Sprint(hostInfo["nickName"])
	title := fmt.Sprint(roomInfo["name"])
	details := fmt.Sprint(hostInfo["signature"])
	img := fmt.Sprint(roomInfo["photo"])
	video := fmt.Sprint(videoInfo["streamurl"])
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = details
	live.LivingIMG = img
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
