package getters

import (
    "errors"
    "strings"
    "github.com/buger/jsonparser"
    "strconv"
    "fmt"
)

//bilibili Bilibili直播
type bilibili struct{}

//Site 实现接口
func (i *bilibili) Site() string { return "Bilibili直播" }

func (i *bilibili) GetExtraInfo(roomid string) (info ExtraInfo, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    infoUrl := "https://api.live.bilibili.com/room/v1/Room/get_info?room_id=" + roomid
    tmp, _ := httpGet(infoUrl)
    info.RoomTitle, _ = jsonparser.GetString([]byte(tmp), "data", "title")
    info.Site = "Bilibili"
    info.OwnerName = ""
    info.RoomID = roomid
    return
}

//SiteURL 实现接口
func (i *bilibili) SiteURL() string {
    return "http://live.bilibili.com"
}

//SiteIcon 实现接口
func (i *bilibili) SiteIcon() string {
    return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *bilibili) FileExt() string {
    return "flv"
}

//NeedFFMpeg 实现接口
func (i *bilibili) NeedFFMpeg() bool {
    return false
}

//GetRoomInfo 实现接口
func (i *bilibili) GetRoomInfo(url string) (id string, live bool, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    urlsplit := strings.Split(url, "/")
    fakeid := urlsplit[len(urlsplit)-1]
    api := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + fakeid
    tmp, err := httpGet(api)
    idInt, err := jsonparser.GetInt([]byte(tmp), "data", "room_id")
    id = strconv.Itoa(int(idInt))
    live_status, err := jsonparser.GetInt([]byte(tmp), "data", "live_status")
    if live_status == 1 {
        live = true
        fmt.Println(id + " live now")
    } else {
        live = false
    }
    return
}

//GetLiveInfo 实现接口
func (i *bilibili) GetLiveInfo(id string) (live LiveInfo, err error) {
    defer func() {
        if recover() != nil {
            err = errors.New("fail get data")
        }
    }()
    live = LiveInfo{RoomID: id}
    url := "https://api.live.bilibili.com/api/playurl?cid=" + id + "&otype=json&quality=0&platform=web"
    tmp, err := httpGet(url)
    videoLinkNum := 1
    jsonparser.ArrayEach([]byte(tmp), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if videoLinkNum == 1 {
            live.VideoURL, _ = jsonparser.GetString(value, "url")
            fmt.Println("get live url")
        }
        videoLinkNum++
    }, "durl")
    if live.VideoURL == "" {
        err = errors.New("fail get data")
    }
    return
}
