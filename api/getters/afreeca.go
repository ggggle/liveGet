package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//afreeca AfreecaTV
type afreeca struct{}

//Site 实现接口
func (i *afreeca) Site() string { return "AfreecaTV" }

func (i *afreeca) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *afreeca) SiteURL() string {
	return "http://www.afreecatv.com"
}

//SiteIcon 实现接口
func (i *afreeca) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *afreeca) FileExt() string {
	return "ts"
}

//NeedFFMpeg 实现接口
func (i *afreeca) NeedFFMpeg() bool {
	return true
}

//GetRoomInfo 实现接口
func (i *afreeca) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("play\\.afreecatv\\.com/(\\w+)(/\\d+)*")
	id = reg.FindStringSubmatch(url)[1]
	tmp, err := httpGet(url)
	if !strings.Contains(tmp, fmt.Sprintf("szBjId   = '%s'", id)) {
		id = ""
	} else {
		live = strings.Contains(tmp, "\"og:title\" content=\"[생]")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *afreeca) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://play.afreecatv.com/" + id
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("nBroadNo = (\\d+)")
	rid := reg.FindStringSubmatch(tmp)[1]
	tmp, err = httpPost("http://live.afreecatv.com:8057/afreeca/player_live_api.php", "bno="+rid)
	json := *(pruseJSON(tmp).JToken("CHANNEL"))
	nick := fmt.Sprint(json["BJNICK"])
	title := fmt.Sprint(json["TITLE"])
	stpt := fmt.Sprint(json["STPT"])
	img := fmt.Sprintf("http://liveimg.afreecatv.com/%s.gif", rid)
	if stpt == "RTMP" {
		url = fmt.Sprintf("http://sessionmanager01.afreeca.tv:6060/broad_stream_assign.html?return_type=gs_cdn&broad_key=%s-flash-hd-rtmp", rid)
	} else {
		url = fmt.Sprintf("http://resourcemanager.afreeca.tv:9090/broad_stream_assign.html?return_type=gs_cdn&broad_key=%s-flash-hd-hls", rid)
	}
	tmp, err = httpGet(url)
	json = *pruseJSON(tmp)
	video := fmt.Sprint(json["view_url"])
	if stpt == "HLS" {
		tmp, err = httpPost("http://live.afreecatv.com:8057/afreeca/player_live_api.php", "type=pwd&bno="+rid)
		json = *(pruseJSON(tmp).JToken("CHANNEL"))
		aid := fmt.Sprint(json["AID"])
		video += "?aid=" + aid
	} else {
		reg, _ := regexp.Compile("rtmp://g7\\.\\w+")
		if reg.FindString(video) != "" {
			reg, _ := regexp.Compile("rtmp://([\\w\\.]+)/(\\S+)")
			l := reg.FindStringSubmatch(video)
			host := l[1]
			path := l[2]
			host = strings.Join(strings.Split(host, ".")[2:], ".")
			host = strings.Split(path, "/")[0] + "." + host
			video = fmt.Sprintf("rtmp://%s/%s", host, path)
		}
	}
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = ""
	live.LivingIMG = img
	live.VideoURL = video

	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
