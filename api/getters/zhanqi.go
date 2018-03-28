package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//zhanqi 战旗直播
type zhanqi struct{}

//Site 实现接口
func (i *zhanqi) Site() string { return "战旗直播" }

func (i *zhanqi) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *zhanqi) SiteURL() string {
	return "http://www.zhanqi.tv"
}

//SiteIcon 实现接口
func (i *zhanqi) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *zhanqi) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *zhanqi) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *zhanqi) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	html, err := httpGet(url)
	if !strings.Contains(html, "<title>战旗直播_高清流畅的游戏直6播平台 - zhanqi.tv</title>") {
		live = strings.Contains(html, "\"Status\":4")
		reg, _ := regexp.Compile("\"Status\":\\d,\"RoomId\":\\d+")
		tmp := reg.FindString(html)
		reg, _ = regexp.Compile("(\"RoomId\"):(\\d+)")
		id = reg.FindStringSubmatch(tmp)[2]
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *zhanqi) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://apis.zhanqi.tv/static/v2.1/room/"
	url = fmt.Sprintf("%s%s.json", url, id)
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("data"))
	nick := json["nickname"].(string)
	title := json["title"].(string)
	img := json["bpic"].(string)
	key := json["videoId"].(string)
	video := fmt.Sprintf("http://wshdl.cdn.zhanqi.tv/zqlive/%s.flv", key)
	live.LiveNick = nick
	live.RoomTitle = title
	live.LivingIMG = img
	live.RoomDetails = ""
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
