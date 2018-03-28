package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//qie 企鹅直播
type qie struct{}

//Site 实现接口
func (i *qie) Site() string { return "企鹅直播" }

func (i *qie) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *qie) SiteURL() string {
	return "http://live.qq.com"
}

//SiteIcon 实现接口
func (i *qie) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *qie) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *qie) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *qie) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	html, err := httpGet(url)
	reg, _ := regexp.Compile("\"room_id\":\\d+")
	tmp := reg.FindString(html)
	live = strings.Contains(html, "\"show_status\":1,")
	reg, _ = regexp.Compile("\\d+")
	id = reg.FindString(tmp)
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *qie) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.qie.tv/api/v1/"
	args := fmt.Sprintf("room/%s?aid=wp&client_sys=wp&time=%d", id, getUnixTimesTamp())
	url = fmt.Sprintf("%s%s&auth=%s", url, args, getMD5String(args+"zNzMV1y4EMxOHS6I5WKm"))
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("data"))
	video := fmt.Sprintf("%s/%s", json["rtmp_url"], json["rtmp_live"])
	img := json["room_src"].(string)
	title := json["room_name"].(string)
	details := json["show_details"].(string)
	nick := json["nickname"].(string)
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
