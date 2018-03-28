package getters

import (
	"errors"
	"regexp"
	"strings"
)

//qiedianjing 企鹅电竞
type qiedianjing struct{}

//Site 实现接口
func (i *qiedianjing) Site() string { return "企鹅电竞" }

func (i *qiedianjing) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *qiedianjing) SiteURL() string {
	return "http://egame.qq.com"
}

//SiteIcon 实现接口
func (i *qiedianjing) SiteIcon() string {
	return "http://egame.gtimg.cn/club/pgg_pcweb/img/favicons/favicon.ico"
}

//FileExt 实现接口
func (i *qiedianjing) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *qiedianjing) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *qiedianjing) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	reg, _ := regexp.Compile("egame\\.qq\\.com/live\\?anchorid=(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	if id != "" {
		url = "http://share.egame.qq.com/cgi-bin/pgg_skey_async_fcgi?param={%220%22:{%22module%22:%22pgg_live_read_svr%22,%22method%22:%22get_live_and_profile_info%22,%22param%22:{%22anchor_id%22:" + id + "}}}"
		var tmp string
		tmp, err = httpGet(url)
		if strings.Contains(tmp, "\"retMsg\":\"ok\"") && strings.Contains(tmp, "\"provider\": 2") {
			json := pruseJSON(tmp)
			live = (*json.JToken("data").JToken("0").JToken("retBody").JToken("data").JToken("profile_info"))["is_live"].(float64) == 1
		} else {
			id = ""
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *qiedianjing) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://share.egame.qq.com/cgi-bin/pgg_skey_async_fcgi?param={%220%22:{%22module%22:%22pgg_live_read_svr%22,%22method%22:%22get_live_and_profile_info%22,%22param%22:{%22anchor_id%22:" + id + "}}}"
	tmp, err := httpGet(url)
	json := pruseJSON(tmp).JToken("data").JToken("0").JToken("retBody").JToken("data")
	profile_info, video_info := *(json.JToken("profile_info")), *(json.JToken("video_info"))
	nick := profile_info["nick_name"].(string)
	details := profile_info["brief"].(string)
	title := video_info["title"].(string)
	video := (*video_info.JTokens("stream_infos")[0])["play_url"].(string)
	img := video_info["url"].(string)
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
