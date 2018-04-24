package getters

import (
	"errors"
	"regexp"
	"strings"
    "fmt"
    "strconv"
    "encoding/hex"
)

//huya 虎牙直播
type huya struct{}

//SiteURL 实现接口
func (i *huya) SiteURL() string {
	return "http://www.huya.com"
}

//Site 实现接口
func (i *huya) Site() string { return "虎牙直播" }

func (i *huya) GetExtraInfo(string) (info ExtraInfo, err error) { return }

//SiteIcon 实现接口
func (i *huya) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *huya) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *huya) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *huya) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("huya\\.com/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://m.huya.com/" + id
	html, err := httpGetWithUA(url, ipadUA)
	if !strings.Contains(html, "找不到此页面") {
		live = strings.Contains(html, "ISLIVE = true")
	} else {
		err = errors.New("fail get id")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huya) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.huya.com/" + id
	html, err := httpGet(url)
	reg, _ := regexp.Compile("\"channel\":\"*(\\d+)\"*,")
	channel, _ := strconv.Atoi(reg.FindStringSubmatch(html)[1])
	reg, _ = regexp.Compile("\"sid\":\"*(\\d+)\"*,")
	sid, _ := strconv.Atoi(reg.FindStringSubmatch(html)[1])
	_hex := fmt.Sprintf("0000009E10032C3C4C56066C6976657569660D6765744C6976696E67496E666F7D0000750800010604745265711D0000680A0A0300000000000000001620FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF2600361777656226323031377633322E313131302E33266875796146005C0B1300000000%08x2300000000%08x3300000000000000000B8C980CA80C",
	    channel, sid)
	_hex = strings.ToUpper(_hex)
	req_data, _ := hex.DecodeString(_hex)
	video_content, err := httpPost("http://cdn.wup.huya.com/", string(req_data))
	print(video_content)

	reg, _ = regexp.Compile(fmt.Sprintf("(%d-%d[^f]+)", channel, sid))
	video_id := reg.FindStringSubmatch(video_content)[1]

	reg, _ = regexp.Compile("wsSecret=([0-9a-z]{32})")
	wsSecret := reg.FindStringSubmatch(video_content)[1]

	reg, _ = regexp.Compile("wsTime=([0-9a-z]{8})")
	wsTime := reg.FindStringSubmatch(video_content)[1]

	reg, _ = regexp.Compile("://(.+\\.(flv|stream)\\.huya\\.com/(hqlive|huyalive))")
	video_url := reg.FindStringSubmatch(video_content)[1]

	live.VideoURL = fmt.Sprintf("http://%s/%s.flv?wsSecret=%s&wsTime=%s", video_url, video_id, wsSecret, wsTime)

	url = "http://m.huya.com/" + id
	html, err = httpGetWithUA(url, ipadUA)
	live.LiveNick = getValue(html, "ANTHOR_NICK")
	live.LivingIMG = getValue(html, "picURL")
	live.RoomDetails = ""
	live.RoomTitle = getValue(html, "liveRoomName")
	return
}

func getValue(data, name string) string {
	reg, _ := regexp.Compile("var " + name + " = [\"'](.*)[\"']")
	return reg.FindStringSubmatch(data)[1]
}
