package getters

import (
    "errors"
    "fmt"
    "strings"
    "github.com/buger/jsonparser"
)

//douyu 斗鱼直播
type douyu struct{}

//Site 实现接口
func (i *douyu) Site() string { return "斗鱼直播" }

func (i *douyu) GetExtraInfo(room string) (info ExtraInfo, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    info.Site = "[斗鱼]"
    url := "http://open.douyucdn.cn/api/RoomApi/room/" + room
    json, _ := httpGet(url)
    if len(json) > 0 {
        errorNo, _ := jsonparser.GetInt([]byte(json), "error")
        if 0 != errorNo {
            errinfo, _ := jsonparser.GetString([]byte(json), "data")
            err = errors.New("json错误码为:" + string(errorNo) + " info:" + errinfo)
            return
        } else {
            info.RoomTitle, _ = jsonparser.GetString([]byte(json), "data", "room_name")
            info.RoomID, _ = jsonparser.GetString([]byte(json), "data", "room_id")
            info.CateName, _ = jsonparser.GetString([]byte(json), "data", "cate_name")
            info.RoomStatus, _ = jsonparser.GetString([]byte(json), "data", "room_status")
            info.StartTime, _ = jsonparser.GetString([]byte(json), "data", "start_time")
            info.OwnerName, _ = jsonparser.GetString([]byte(json), "data", "owner_name")
        }
    }
    return
}

//SiteURL 实现接口
func (i *douyu) SiteURL() string {
    return "http://www.douyu.com"
}

//SiteIcon 实现接口
func (i *douyu) SiteIcon() string {
    return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *douyu) FileExt() string {
    return "flv"
}

//NeedFFMpeg 实现接口
func (i *douyu) NeedFFMpeg() bool {
    return false
}

//GetRoomInfo 实现接口
func (i *douyu) GetRoomInfo(url string) (id string, live bool, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    urlSplit := strings.Split(url, "/")
    room := urlSplit[len(urlSplit)-1]
    url = "http://open.douyucdn.cn/api/RoomApi/room/" + room
    json, _ := httpGet(url)
    if len(json) > 0 {
        errorNo, _ := jsonparser.GetInt([]byte(json), "error")
        if 0 != errorNo {
            errinfo, _ := jsonparser.GetString([]byte(json), "data")
            err = errors.New("json错误码为:" + string(errorNo) + " info:" + errinfo)
            return
        } else
        {
            id, _ = jsonparser.GetString([]byte(json), "data", "room_id")
            room_status, _ := jsonparser.GetString([]byte(json), "data", "room_status")
            //1直播中  2未开播
            if 0 == strings.Compare(room_status, "1") {
                live = true
            } else {
                live = false
            }
        }
    }
    return
}

//GetLiveInfo 实现接口
func (i *douyu) GetLiveInfo(id string) (live LiveInfo, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data panic")
        }
    }()
    live = LiveInfo{RoomID: id}
    url := "http://www.douyutv.com/api/v1/"
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
        err = errors.New("fail get data no video")
    }
    return
}
