package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"strings"
	"github.com/ggggle/luzhibo/api"
	"github.com/ggggle/luzhibo/workers"
	"github.com/ggggle/luzhibo/api/getters"
	"os/exec"
)

var tasks []task

type task struct {
	API       *api.LuzhiboAPI
	worker    workers.Worker
	startTime time.Time
}

func addTask(oa *api.LuzhiboAPI, path string, t, s bool) bool {
	var w workers.Worker
	var e error
	if oa.NeedFFmpeg && !hasFFmpeg() {
		return false
	}
	if !t {
		w, e = workers.NewSingleWorker(oa, path, nil)
	} else {
		w, e = workers.NewMultipleWorker(oa, path, nil)
	}
	if e != nil {
		return false
	}
	tt := task{API: oa, worker: w}
	tasks = append(tasks, tt)
	if s {
		startTask(len(tasks) - 1)
	}
	return true
}

func addTaskEx(url, path string, t, s bool) bool {
	oa := api.New(url)
	if oa == nil {
		return false
	}
	return addTask(oa, path, t, s)
}

func init() {
	tasks = make([]task, 0)
}

func addTasks(urls string) int {
	c := 0
	list := strings.Split(urls, "\n")
	ch := make(chan bool, 0)
	i := 0
	l := len(list)
	for _, url := range list {
		go func(u string) {
			oa := api.New(u)
			if oa != nil {
				i, _, e := oa.GetRoomInfo()
				if e == nil {
					p := fmt.Sprintf("[%s]%s_%s", oa.Site, i, time.Now().Format("20060102150405"))
					if (addTask(oa, p, true, true)) {
						c++
					}
				}
			}
			i++
			if i == l {
				ch <- true
			}
		}(url)
	}
	<-ch
	return c
}

func delTask(i int, f bool) bool {
	if i >= 0 && i+1 <= len(tasks) {
		w := tasks[i].worker
		w.Stop()
		if f {
			_, _, _, p, _ := w.GetTaskInfo(false)
			delPath(p)
		}
		a, b := tasks[:i], tasks[i+1:]
		r := append(a, b...)
		tasks = r
		return true
	}
	return false
}

func startTask(x int) bool {
	if x >= 0 && x+1 <= len(tasks) {
		v := tasks[x]
		w, e := v.worker.Restart()
		if e == nil {
			tasks[x].worker = w
			tasks[x].startTime = time.Now().Local()
			return true
		}
	}
	return false
}

func stopTask(x int) bool {
	if x >= 0 && x+1 <= len(tasks) {
		tasks[x].worker.Stop()
		return true
	}
	return false
}

func updateTaskStatus(x int) {
	if x >= 0 && x+1 <= len(tasks) {
		v := tasks[x].worker
		if _, r, _, _, _ := v.GetTaskInfo(false); !r {
			startTask(x)
		} else {
			v.Stop()
		}
	}
}

type taskInfo struct {
	Site       string
	SiteIcon   string
	SiteURL    string
	URL        string
	ID         string
	Live       bool
	M          bool
	Run        bool
	FileExt    string
	NeedFFmpeg bool
	Files      []string
	Path       string
	Index      int64
	StartTime  string
	TimeLong   string
	LiveInfo   *getters.LiveInfo
}

func getTaskInfo(x int) (o *taskInfo, te int) {
	if x < 0 || x+1 > len(tasks) {
		te = -4
		return
	}
	o = &taskInfo{}
	if x+1 <= len(tasks) {
		v := tasks[x]
		o.Site = v.API.Site
		o.SiteURL = v.API.SiteURL
		o.SiteIcon = v.API.Icon
		o.URL = v.API.URL
		o.NeedFFmpeg = v.API.NeedFFmpeg
		o.FileExt = v.API.FileExt
		i, l, e := v.API.GetRoomInfo()
		if e == nil {
			o.ID = i
			o.Live = l
		} else {
			te--
		}
		tt, r, ind, p, inf := v.worker.GetTaskInfo(true)
		o.M = tt == 2
		o.Run = r
		o.Path = p
		o.Files = getFiles(p)
		o.Index = ind
		o.LiveInfo = inf
		if inf == nil {
			te -= 2
		}
		o.StartTime = v.startTime.Format("2006-01-02:15:04:05")
		if r {
			o.TimeLong = timeLongToStr(time.Now().Local().Sub(v.startTime))
		}
	}
	return
}

func timeLongToStr(v time.Duration) string {
	cm := int64(60)
	ch := int64(cm * 60)
	cd := int64(ch * 24)
	ts := int64(v.Seconds())
	d := ts / cd
	h := ts % cd / ch
	m := ts % cd % ch / cm
	s := ts % cd % ch % cm
	str := fmt.Sprintf("%d天%02d小时%02d分钟%02d秒", d, h, m, s)
	return str
}

func getTaskInfoList() (list []*taskInfo, err, e bool) {
	l := len(tasks)
	if l == 0 {
		e = true
	} else {
		t := make([]*taskInfo, l)
		ch := make(chan bool, 0)
		c := 0
		for i := 0; i < l; i++ {
			go func(x int) {
				v, te := getTaskInfo(x)
				if te != 0 {
					err = true
				}
				if v != nil {
					t[x] = v
				}
				c++
				if c == l {
					ch <- true
				}
			}(i)
		}
		<-ch
		list = t
	}
	return
}

func pathExist(path string) (bool, bool) {
	p, err := os.Stat(path)
	if err == nil {
		return true, p.IsDir()
	}
	return os.IsExist(err), false
}

func getFiles(path string) []string {
	l := make([]string, 0)
	if e, d := pathExist(path); e {
		if d {
			files, err := ioutil.ReadDir(path)
			if err == nil {
				for _, f := range files {
					if !f.IsDir() {
						p := path + "/" + f.Name()
						l = append(l, p)
					}
				}
			} else {
				l = nil
			}
		} else {
			l = append(l, path)
		}
		return l
	}
	return nil
}

func delPath(path string) {
	os.RemoveAll(path)
}

func hasFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}
