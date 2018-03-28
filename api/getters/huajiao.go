package getters

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	json2 "github.com/Baozisoftware/golibraries/json"
)

//huajiao 花椒直播
type huajiao struct{}

//Site 实现接口
func (i *huajiao) Site() string { return "花椒直播" }

func (i *huajiao) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *huajiao) SiteURL() string {
	return "http://www.huajiao.com"
}

//SiteIcon 实现接口
func (i *huajiao) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *huajiao) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *huajiao) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *huajiao) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	reg, _ := regexp.Compile("huajiao\\.com/user/(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://webh.huajiao.com/User/getUserFeeds?uid=" + id
	html, err := httpGet(url)
	if err == nil {
		if strings.Contains(html, "\"data\":[]}") {
			id = ""
		} else {
			live = strings.Contains(html, "\"replay_status\":0")
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huajiao) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.huajiao.com/user/" + id
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("\"nickname\":(\"[^\"]*\")")
	nick := reg.FindStringSubmatch(tmp)[1]
	nick, _ = strconv.Unquote(nick)
	url = "http://webh.huajiao.com/User/getUserFeeds?uid=" + id
	tmp, err = httpGet(url)
	json := pruseJSON(tmp).JToken("data")
	feeds := json.JTokens("feeds")
	var tf *json2.JObject
	for _, v := range feeds {
		if (*v)["type"].(float64) == 1 {
			tf = v.JToken("feed")
			break
		}
	}
	if tf == nil {
		err = errors.New("fail get data")
		return
	}
	feed := *tf
	sn := feed["sn"]
	img := feed["image"].(string)
	title := feed["title"].(string)
	url = fmt.Sprintf("http://g2.live.360.cn/liveplay?stype=flv&channel=live_huajiao_v2&bid=huajiao&sn=%s&sid=null&_rate=null&ts=null", sn)
	tmp, err = httpGet(url)
	tmp = tmp[0:3] + tmp[6:]
	bytes, err := base64.StdEncoding.DecodeString(tmp)
	tmp = string(bytes)
	json = pruseJSON(tmp)
	video := (*json)["main"].(string)
	live.LiveNick = nick
	live.LivingIMG = img
	live.RoomDetails = ""
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
