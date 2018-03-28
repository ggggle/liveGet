package getters

import (
	"errors"
	"regexp"
	"strings"
	"fmt"
)

//inke 映客直播
type inke struct{}

//Site 实现接口
func (i *inke) Site() string { return "映客直播" }

func (i *inke) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *inke) SiteURL() string {
	return "http://www.inke.cn"
}

//SiteIcon 实现接口
func (i *inke) SiteIcon() string {
	return "http://static.inke.com/s/images/favicon.ico"
}

//FileExt 实现接口
func (i *inke) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *inke) NeedFFMpeg() bool {
	return true
}

//GetRoomInfo 实现接口
func (i *inke) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("inke\\.cn/live.html\\?uid=\\d+&id=(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://webapi.busi.inke.cn/web/live_share_pc?uid=0&id=" + id
	tmp, err := httpGet(url)
	if strings.Contains(tmp, "\"inke_id\":0") {
		id = ""
	} else {
		live = strings.Contains(tmp, "\"status\":1")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *inke) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://webapi.busi.inke.cn/web/live_share_pc?uid=0&id=" + id
	tmp, err := httpGet(url)
	json := pruseJSON(tmp).JToken("data")
	info1, info2 := *(json.JToken("media_info")), *(json.JToken("file"))
	nick := fmt.Sprint(info1["nick"])
	title := fmt.Sprint(info2["title"])
	details := fmt.Sprint(info1["description"])
	img := fmt.Sprint(info2["pic"])
	video := fmt.Sprint(info2["record_url"])
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
