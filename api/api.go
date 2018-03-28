package api

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ggggle/luzhibo/api/getters"
)

//LuzhiboAPI API object
type LuzhiboAPI struct {
	Id         string
	URL        string
	G          getters.Getter
	Site       string
	SiteURL    string
	Icon       string
	FileExt    string
	NeedFFmpeg bool
}

//New 使用网址创建一个实例
func New(url string) *LuzhiboAPI {
	var r *LuzhiboAPI
	g := getGetter(url)
	if g != nil {
		i := &LuzhiboAPI{}
		i.G = g
		i.URL = url
		i.Site = g.Site()
		i.SiteURL = g.SiteURL()
		i.Icon = g.SiteIcon()
		i.FileExt = g.FileExt()
		i.NeedFFmpeg = g.NeedFFMpeg()
		r = i
	} else {
		r = nil
	}
	return r
}

//GetRoomInfo 取直播间信息
func (i *LuzhiboAPI) GetRoomInfo() (id string, live bool, err error) {
	if i.URL == "" || i.G == nil {
		err = errors.New("not has url or not found getter")
		return
	}
	id, live, err = i.G.GetRoomInfo(i.URL)
	i.Id = id
	if Logger != nil {
		s := fmt.Sprintf("获取房间信息\"%s\",结果:", i.URL)
		if err == nil {
			s += fmt.Sprintf("成功(直播平台:\"%s\",房间ID:\"%s\",已开播:", i.Site, id)
			if live {
				s += fmt.Sprint("\"是\".).")
			} else {
				s += fmt.Sprint("\"否\".).")
			}
		} else {
			s += fmt.Sprint("失败(获取时出错).")
			s += fmt.Sprint(err.Error())
			Logger.Print(s)
		}
	}
	return
}

//GetLiveInfo 取直播信息
func (i *LuzhiboAPI) GetLiveInfo() (live getters.LiveInfo, err error) {
	if i.Id == "" || i.G == nil {
		err = errors.New("not has id or not found getter")
		return
	}
	live, err = i.G.GetLiveInfo(i.Id)
	if Logger != nil {
		s := fmt.Sprintf("获取直播信息\"%s\",结果:", i.URL)
		if err == nil {
			s += fmt.Sprintf("成功(直播平台:\"%s\",房间ID:\"%s\",房间标题:\"%s\",主播昵称:\"%s\",直播流地址:\"%s\".).", i.Site, i.Id, live.RoomTitle, live.LiveNick, live.VideoURL)
		} else {
			s += fmt.Sprint("失败(获取时出错).")
			s += fmt.Sprint(err.Error())
		}
		Logger.Print(s)
	}
	return
}

func getGetter(url string) getters.Getter {
	url = strings.ToLower(url)
	regs := []string{"(douyu\\.tv)|((douyu)|(douyutv)\\.com)",
		"www\\.panda\\.tv",
		"zhanqi\\.tv",
		"longzhu\\.com",
		"huya\\.com",
		"live\\.qq\\.com",
		"live\\.bilibili\\.com",
		"quanmin\\.tv",
		"huajiao\\.com",
		"huomao\\.com",
		"yizhibo\\.com",
		"egame.qq\\.com",
		"chushou\\.tv",
		"inke\\.cn",
		"play\\.afreecatv\\.com",
		"xingyan\\.panda\\.tv"}
	gs := getters.Getters()
	for i := 0; i < len(gs); i++ {
		if ok, _ := regexp.MatchString(regs[i], url); ok {
			return gs[i]
		}
	}
	return nil
}

func GetSupports() []string {
	gs := getters.Getters()
	ret := make([]string, len(gs))
	for i, oa := range gs {
		ret[i] = oa.Site()
		if oa.NeedFFMpeg() {
			ret[i] = ret[i] + "*"
		}
	}
	return ret
}

var Logger *log.Logger
