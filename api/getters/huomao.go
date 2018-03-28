package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//huomao 火猫直播
type huomao struct{}

//Site 实现接口
func (i *huomao) Site() string { return "火猫直播" }

func (i *huomao) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteURL 实现接口
func (i *huomao) SiteURL() string {
	return "http://www.huomao.com"
}

//SiteIcon 实现接口
func (i *huomao) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *huomao) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *huomao) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *huomao) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("var cid = (\\d+);")
	id = reg.FindStringSubmatch(tmp)[1]
	url = getURL(id)
	tmp, err = httpGet(url)
	if strings.Contains(tmp, "\"is_live\"") {
		live = strings.Contains(tmp, "\"is_live\":1")
	} else {
		err = errors.New("fail get data")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huomao) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := getURL(id)
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("data"))
	nick := json["username"].(string)
	title := json["channel"].(string)
	details := json["content"].(string)
	video := (*(json.JTokens("streamList"))[2])["HD"].(string)
	live.LiveNick = nick
	live.LivingIMG = ""
	live.RoomDetails = details
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}

func getURL(id string) string {
	time := getUnixTimesTamp()
	args := fmt.Sprintf("%d%sEU*T*)*(#23ssdfd", time, id)
	token := getMD5String(args)
	url := fmt.Sprintf("http://api.huomao.com/channels/channelDetail?&cid=%s&time=%d&token=%s", id, time, token)
	return url
}
