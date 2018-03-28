package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
    "github.com/buger/jsonparser"
)

//longzhu 龙珠直播
type longzhu struct{}

//Site 实现接口
func (i *longzhu) Site() string { return "龙珠直播" }

func (i *longzhu) GetExtraInfo(roomid string) (info ExtraInfo, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    info.Site = "[龙珠]"
    info.RoomID = roomid
    url := "http://roomapicdn.plu.cn/room/RoomAppStatusV2?domain=" + roomid
    json, _ := httpGet(url)
    if len(json) > 0 {
        info.RoomTitle, _ = jsonparser.GetString([]byte(json), "BaseRoomInfo", "BoardCastTitle")
        info.OwnerName, _ = jsonparser.GetString([]byte(json), "BaseRoomInfo", "Name")
    }
    return
}

//SiteURL 实现接口
func (i *longzhu) SiteURL() string {
	return "http://www.longzhu.com"
}

//SiteIcon 实现接口
func (i *longzhu) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *longzhu) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *longzhu) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *longzhu) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("longzhu\\.com/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	if id != "" {
		url := "http://roomapicdn.plu.cn/room/RoomAppStatusV2?domain=" + id
		var tmp string
		tmp, err = httpGet(url)
		if err == nil {
			if strings.Contains(tmp, "IsBroadcasting") {
				live = strings.Contains(tmp, "\"IsBroadcasting\":true")
			} else {
				id = ""
			}
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *longzhu) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{}
	url := "http://roomapicdn.plu.cn/room/RoomAppStatusV2?domain=" + id
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("BaseRoomInfo"))
	nick := json["Name"].(string)
	title := json["BoardCastTitle"].(string)
	details := json["Desc"].(string)
	_id := json["Id"]
	live.RoomID = fmt.Sprintf("%.f", _id)
	url = "http://livestream.plu.cn/live/getlivePlayurl?roomId=" + live.RoomID
	tmp, err = httpGet(url)
	json = *(pruseJSON(tmp).JTokens("playLines")[0].JTokens("urls")[0])
	video := json["securityUrl"].(string)
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = details
	live.LivingIMG = ""
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
