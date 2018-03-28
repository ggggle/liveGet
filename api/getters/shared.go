package getters

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	nhttp "github.com/Baozisoftware/golibraries/http"
	"github.com/Baozisoftware/golibraries/json"
)

//实现一些通用函数/结构

func httpGetWithUA(url, ua string) (data string, err error) {
	resp, err := httpGetResp(url, ua)
	if err == nil {
		var body []byte
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			data = string(body)
		}
	}
	return
}

func httpGet(url string) (data string, err error) {
	return httpGetWithUA(url, "")
}

func httpGetResp(url, ua string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		client := nhttp.NewHttpClient()
		client.SetResponseHeaderTimeout(30)
		client.SetProxy(Proxy)
		req.Header.Set("User-Agent", ua)
		resp, err = client.Do(req)
	}
	return
}

func httpPostResp(url, ua, data string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	if err == nil {
		client := nhttp.NewHttpClient()
		client.SetProxy(Proxy)
		req.Header.Set("User-Agent", ua)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err = client.Do(req)
	}
	return
}

func httpPostWithUA(url, ua, data string) (result string, err error) {
	resp, err := httpPostResp(url, ua, data)
	if err == nil {
		var body []byte
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			result = string(body)
		}
	}
	return
}

func httpPost(url, data string) (result string, err error) {
	return httpPostWithUA(url, "", data)
}

func getUnixTimesTamp() int64 {
	return time.Now().Unix()
}

func getMD5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	r := m.Sum(nil)
	return hex.EncodeToString(r)
}

func forEach(list interface{}, f func(interface{}) bool, maxCount int) (result []interface{}, err error) {
	t := 0
	switch reflect.TypeOf(list).Kind() {
	case reflect.Array:
		t = 1
	case reflect.Slice:
		t = 2
	default:
		err = errors.New("list type error")
	}
	if maxCount < 0 {
		err = errors.New("list count error")
	} else if t > 1 {
		maxCount--
		value := reflect.ValueOf(list)
		count := value.Len()
		tmp := make([]interface{}, 0)
		defer func() {
			if recover() != nil {
				err = errors.New("fild to for each")
			}
		}()
		for i := 0; i < count; i++ {
			if v := value.Index(i).Interface(); f(v) {
				tmp = append(tmp, v)
			}
			if i == maxCount {
				break
			}
		}
		if len(tmp) > 0 {
			result = tmp
		} else {
			err = errors.New("not find")
		}
	}
	return
}

func forEachOne(list interface{}, f func(interface{}) bool) (result interface{}, err error) {
	tmp, err := forEach(list, f, 1)
	if err == nil {
		result = tmp[0]
	}
	return
}

func pruseJSON(data string) *json.JObject {
	return json.PruseJSON(data)
}

//LiveInfo 直播间信息结构
type LiveInfo struct {
	RoomTitle   string
	LivingIMG   string
	VideoURL    string
	RoomDetails string
	RoomID      string
	LiveNick    string
}

//获取一些额外信息
type ExtraInfo struct {
    Site         string   //网站
    RoomTitle    string   //房间标题
    RoomID       string   //房间id
    CateName     string   //分类
    RoomStatus   string   //直播状态
    StartTime    string   //开始时间
    OwnerName    string   //主播名
}

//Getter 房间/直播信息获取接口
type Getter interface {
    GetExtraInfo(string) (ExtraInfo, error)   //从其它接口获取一些需要的信息
	GetRoomInfo(string) (string, bool, error) //获取房间信息,参数为房间地址,返回房间号,是否开播
	GetLiveInfo(string) (LiveInfo, error)     //获取直播信息,参数为房间号,返回直播信息
	Site() string                             //返回平台名称
	SiteURL() string                          //返回平台首页
	SiteIcon() string                         //网站图标
	NeedFFMpeg() bool                         //是否需要FFmpeg
	FileExt() string                          //文件扩展名
}

//Getters 所有获取接口
func Getters() []Getter {
	return []Getter{&douyu{}, &panda{}, &zhanqi{}, &longzhu{}, &huya{}, &qie{}, &bilibili{}, &quanmin{}, &huajiao{}, &huomao{}, &yi{}, &qiedianjing{}, &chushou{}, &inke{}, &afreeca{}, &xingyan{}}
}

var Proxy = ""

const ipadUA = "Mozilla/5.0 (iPad; CPU OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4"
