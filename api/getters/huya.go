package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//huya 虎牙直播
type huya struct{}

//SiteURL 实现接口
func (i *huya) SiteURL() string {
	return "http://www.huya.com"
}

//Site 实现接口
func (i *huya) Site() string { return "虎牙直播" }

func (i *huya) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteIcon 实现接口
func (i *huya) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *huya) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *huya) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *huya) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("huya\\.com/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://m.huya.com/" + id
	html, err := httpGetWithUA(url, ipadUA)
	if !strings.Contains(html, "找不到此页面") {
		live = strings.Contains(html, "ISLIVE = true")
	} else {
		err = errors.New("fail get id")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huya) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://m.huya.com/" + id
	html, err := httpGetWithUA(url, ipadUA)
	title := getValue(html, "liveRoomName")
	nick := getValue(html, "ANTHOR_NICK")
	img := getValue(html, "picURL")
	reg, _ := regexp.Compile("cid: '(\\d+/\\d+)'")
	cid := strings.Replace(reg.FindStringSubmatch(html)[1], "/", "_", 1)
	video := fmt.Sprintf("http://hls.yy.com/%s_100571200.flv", cid)
	live.LiveNick = nick
	live.LivingIMG = img
	live.RoomDetails = ""
	live.RoomTitle = title
	live.VideoURL = video
	if live.VideoURL == "" {
		err = errors.New("fail get data")
	}
	return
}

func getValue(data, name string) string {
	reg, _ := regexp.Compile("var " + name + " = [\"'](.*)[\"']")
	return reg.FindStringSubmatch(data)[1]
}
