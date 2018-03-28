package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Baozisoftware/GoldenDaemon"
	nhttp "github.com/Baozisoftware/golibraries/http"
	"github.com/ggggle/luzhibo/api"
	"github.com/inconshreveable/go-update"
)

type checkRet struct {
	Pass    bool
	Has     bool
	Live    bool
	Err     bool
	Path    string
	FileExt string
	Support bool
}

type tasksRet struct {
	Tasks []*taskInfo
	Err   bool
	E     bool
}

type ajaxHandler struct{}

//ServeHTTP 实现接口
func (_ ajaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	act := r.Form.Get("act")
	switch act {
	case "check":
		tr := checkRet{}
		url := r.Form.Get("url")
		oa := api.New(url)
		if oa == nil {
			tr.Pass = false
		} else {
			tr.Pass = true
			i, l, e := oa.GetRoomInfo()
			if e == nil {
				tr.Has = true
				tr.Path = fmt.Sprintf("[%s]%s_%s", oa.Site, i, time.Now().Format("20060102150405"))
				tr.Live = l
				tr.FileExt = oa.FileExt
				if oa.NeedFFmpeg {
					tr.Support = hasFFmpeg()
				} else {
					tr.Support = true
				}

			} else {
				tr.Err = true
			}
		}
		j, _ := json.Marshal(tr)
		w.Write(j)
		return
	case "add":
		url, m, p, s := r.Form.Get("url"), r.Form.Get("m"), r.Form.Get("path"), r.Form.Get("run")
		url, _ = nurl.QueryUnescape(url)
		p, _ = nurl.QueryUnescape(p)
		mm, ss := m == "true", s == "true"
		if url != "" && p != "" {
			if addTaskEx(url, p, mm, ss) {
				w.Write([]byte("ok"))
				return
			}
		}
	case "addex":
		urls := r.Form.Get("urls")
		urls, _ = nurl.QueryUnescape(urls)
		i := addTasks(urls)
		w.Write([]byte(strconv.Itoa(i)))
	case "del":
		i, d := r.Form.Get("id"), r.Form.Get("f")
		b := d == "true"
		c, e := strconv.Atoi(i)
		if e == nil {
			if delTask(c-1, b) {
				w.Write([]byte("ok"))
				return
			}
		}

	case "start":
		i := r.Form.Get("id")
		if startOrStopTask(i, true) {
			w.Write([]byte("ok"))
			return
		}
	case "stop":
		i := r.Form.Get("id")
		if startOrStopTask(i, false) {
			w.Write([]byte("ok"))
			return
		}
	case "tasks":
		list, o, e := getTaskInfoList()
		r := tasksRet{}
		r.Err = o
		r.Tasks = list
		r.E = e
		j, _ := json.Marshal(r)
		w.Write(j)
		return
	case "exist":
		p := r.Form.Get("path")
		if pp, _ := pathExist(p); pp {
			w.Write([]byte("exist"))
			return
		}
	case "get":
		i, s := r.Form.Get("id"), r.Form.Get("sub")
		ii, e := strconv.Atoi(i)
		if e == nil {
			inf, _ := getTaskInfo(ii - 1)
			fp := inf.Path
			if s != "" {
				fp += "/" + s + "." + inf.FileExt
			}
			pp := inf.Path
			if inf.M {
				if s != "" {
					pp += "_" + s
				}
				pp += "." + inf.FileExt
			}
			w.Header().Add("Content-Disposition", "attachment; filename=\""+nurl.QueryEscape(pp)+"\"")
			w.Header().Add("Content-Type", "video/x-"+inf.FileExt)
			getAct(fp, w)
		}
		return
	case "ver":
		w.Write([]byte(checkUpdate()))
		return
	case "supports":
		lines := api.GetSupports()
		s := strings.Join(lines, "/")
		w.Write([]byte(s))
		return
	case "update":
		b := doUpdate()
		r := "false"
		if b {
			r = "true"
		}
		w.Write([]byte(r))
		return
	case "quit":
		w.Write([]byte("ok"))
		go func() {
			time.Sleep(time.Second)
			quit()
		}()
	case "log":
	    content, _ := ioutil.ReadFile("lzb.log")
		w.Write(content)
		return
	}
	w.Write([]byte(""))
}

type uiHandler struct{}

//ServeHTTP 实现接口
func (_ uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b64 := base64.StdEncoding
	switch r.URL.Path {
	case "/":
		h := ui_main
		r.ParseForm()
		if r.Form.Get("hta") == "true" {
			h = strings.Replace(h, "hta = false", "hta = true", 1)
		}
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(h))
	case "/hta":
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(hta))
	case "/favicon.ico":
		w.Header().Add("Content-Type", "image/x-icon")
		data, _ := b64.DecodeString(favicon_ico)
		w.Write(data)
	case "/bootstrap.min.css":
		w.Header().Add("Content-Type", "text/css")
		data, _ := b64.DecodeString(bootstrap_min_css)
		w.Write(data)
	case "/bootstrap.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(bootstrap_min_js)
		w.Write(data)
	case "/jquery.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(jquery_min_js)
		w.Write(data)
	case "/flv.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(flv_min_js)
		w.Write(data)
	case "/fonts/glyphicons-halflings-regular.woff2":
		w.Header().Add("Content-Type", "application/octet-stream")
		data, _ := b64.DecodeString(glyphicons_halflings_regular_woff2)
		w.Write(data)
	case "/fonts/glyphicons-halflings-regular.eot":
		w.Header().Add("Content-Type", "application/octet-stream")
		data, _ := b64.DecodeString(glyphicons_halflings_regular_eot)
		w.Write(data)
	}
}

func getFile(path string, w http.ResponseWriter) {
	f, e := os.Open(path)
	defer f.Close()
	eof := false
	if e == nil {
		buf := make([]byte, bytes.MinRead)
		for {
			t, e := f.Read(buf)
			if e != nil {
				if e == io.EOF {
					eof = true
				} else {
					break
				}
			}
			_, e = w.Write(buf[:t])
			if e != nil || eof {
				break
			}
		}
	}
}

func getDir(path string, w http.ResponseWriter) {
	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() {
				p := path + "/" + f.Name()
				getFile(p, w)
			}
		}
	}
}

func getAct(path string, w http.ResponseWriter) {
	if pe, d := pathExist(path); pe {
		if d {
			getDir(path, w)
		} else {
			getFile(path, w)
		}
	} else {
		w.Write([]byte("no exist"))
	}

}

func startOrStopTask(i string, m bool) bool {
	c, e := strconv.Atoi(i)
	if e != nil {
		return false
	}
	c--
	if m {
		return startTask(c)
	}
	return stopTask(c)
}

func startServer(s string) {
	http.Handle("/", uiHandler{})
	http.Handle("/bootstrap.min.css", uiHandler{})
	http.Handle("/bootstrap.min.js", uiHandler{})
	http.Handle("/jquery.min.css", uiHandler{})
	http.Handle("/flv.min.css", uiHandler{})
	http.Handle("/ajax", ajaxHandler{})
	http.ListenAndServe(s, nil)
	os.Exit(0)
}

func httpGet(url string) (data string, err error) {
	client := nhttp.NewHttpClient()
	client.SetProxy(*proxy)
	resp, err := client.GetResp(url)
	if err == nil {
		var body []byte
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				data = string(body)
			}
		} else {
			err = errors.New("resp StatusCode is not 200.")
		}
	}
	return
}

func checkUpdate() string {
	s := fmt.Sprint("更新检测,结果:")
	r := strconv.Itoa(ver) + "|"
	if updated || updatting {
		r += "null"
	} else {
		data, err := httpGet("https://api.github.com/repos/Baozisoftware/luzhibo/releases/latest")
		if err == nil {
			if data != "" {
				reg, _ := regexp.Compile("Ver (\\d{10})")
				data = reg.FindStringSubmatch(data)[1]
				if v, _ := strconv.Atoi(data); v > ver {
					s += fmt.Sprintf("有新版本(%d->%d).", ver, v)
					r += data
				} else {
					s += fmt.Sprintf("无新版本(当前版本:%d).", ver)
					r += "null"
				}
			} else {
				s += fmt.Sprint("获取服务器版本号失败.")
				r += "null"
			}
		} else {
			s += fmt.Sprint("获取服务器版本号失败.")
			r += "null"
		}
	}
	logger.Print(s)
	return r
}

var updated = false
var updatting = false

func doUpdate() bool {
	if updatting {
		return false
	}
	if updated {
		return true
	}
	updatting = true
	url := fmt.Sprintf("https://github.com/Baozisoftware/luzhibo/releases/download/latest/luzhibo_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		url += ".exe"
	}
	resp, err := http.Get(url)
	if err != nil {
		updatting = false
		return false
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{OldSavePath: ""})
	if err != nil {
		updatting = false
		return false
	}
	updated = true
	updatting = false
	if len(tasks) == 0 {
		go func() {
			time.Sleep(time.Second)
			restartSelf()
		}()
	}
	return true
}

func restartSelf() {
	if runtime.GOOS == "windows" || *nt {
		if htaproc != nil {
			htaproc.Kill()
		}
		GoldenDaemon.RestartSelf()
	}
}
