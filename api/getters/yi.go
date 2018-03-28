package getters

import (
	"errors"
	"regexp"
	"strings"
)

//yi 一直播
type yi struct{}

//Site 实现接口
func (i *yi) Site() string { return "一直播" }

func (i *yi) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *yi) SiteURL() string {
	return "http://www.yizhibo.com"
}

//SiteIcon 实现接口
func (i *yi) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *yi) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *yi) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *yi) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	reg, _ := regexp.Compile("yizhibo\\.com/member/personel/user_info\\?memberid=(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://www.yizhibo.com/member/personel/user_works?memberid=" + id
	html, err := httpGet(url)
	if err == nil {
		if strings.Contains(html, "window.location=\"/404.html\";") {
			id = ""
		} else {
			live = strings.Contains(html, "index_all_common index_zb")
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *yi) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.yizhibo.com/member/personel/user_works?memberid=" + id
	html, err := httpGet(url)
	if err == nil {
		reg, _ := regexp.Compile("/l/(\\S+)\\.html")
		id = reg.FindStringSubmatch(html)[1]
		if id != "" {
			url = "http://api.xiaoka.tv/live/web/get_play_live?scid=" + id
			html, err = httpGet(url)
			json := *(pruseJSON(html).JToken("data"))
			nick := json["nickname"].(string)
			title := json["title"].(string)
			video := json["linkurl"].(string)
			img := json["cover"].(string)
			img = "http://alcdn.img.xiaoka.tv/" + img
			live.LiveNick = nick
			live.RoomTitle = title
			live.RoomDetails = ""
			live.LivingIMG = img
			live.VideoURL = video
		}
	}
	if live.VideoURL == "" {
		err = errors.New("fail get data")
	}
	return
}
