package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"strconv"
)

//quanmin 全民直播
type quanmin struct{}

//Site 实现接口
func (i *quanmin) Site() string { return "全民直播" }

func (i *quanmin) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *quanmin) SiteURL() string {
	return "http://www.quanmin.tv"
}

//SiteIcon 实现接口
func (i *quanmin) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *quanmin) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *quanmin) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *quanmin) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("quanmin\\.tv/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	url = fmt.Sprintf("http://www.quanmin.tv/json/rooms/%s/noinfo.json", id)
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp))
	live = json["play_status"].(bool)
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *quanmin) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := fmt.Sprintf("http://www.quanmin.tv/json/rooms/%s/noinfo.json", id)
	tmp, err := httpGet(url)
	json := *pruseJSON(tmp)
	nick := json["nick"].(string)
	title := json["title"].(string)
	details := json["intro"].(string)
	img := json["thumb"].(string)
	id=strconv.FormatFloat(json["uid"].(float64),'f',0,64)
	video := fmt.Sprintf("http://flv.quanmin.tv/live/%s.flv", id)
	live.LiveNick = nick
	live.LivingIMG = img
	live.RoomDetails = details
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
