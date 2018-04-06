//+build !windows

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ggggle/liveGet/api"
	//"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/pkg/browser"
)

const title = p + " - 控制台"

func cmd() {
	setConsoleTitle()
	t := fmt.Sprintf("---%s (Ver %d)---", p, ver)
	fmt.Println(t)
	for {
		fmt.Println("输入数字并回车确认:")
		fmt.Println("1.添加一个普通任务")
		fmt.Println("2.添加一个循环任务")
		fmt.Println("3.查看当前任务列表")
		fmt.Println("4.打开WebUI(仅部分平台适用)")
		fmt.Println("5.退出程序")
		fmt.Print("请输入:")
		var o int64
		fmt.Scanf("%d\n", &o)
		switch o {
		case 1:
			add(false)
		case 2:
			add(true)
		case 3:
			show()
		case 4:
			openWebUI(!*nhta)
		case 5:
			return
		default:
			fmt.Println("输入错误,请重试!")
		}
	}
}

func show() {
l1:
	fmt.Print("#\t类型\t运行中\t当前序数\t路径\t房间标题\n")
	cc := len(tasks)
	for i := 0; i < cc; i++ {
		o, _ := getTaskInfo(i)
		var st, sr string
		if !o.M {
			st = "普通"
		} else {
			st = "循环"
		}
		if o.Run {
			sr = "运行"
		} else {
			sr = "停止"
		}
		tt := "[无数据]"
		if o.LiveInfo != nil {
			tt = o.LiveInfo.RoomTitle
		}
		fmt.Printf("%d\t%s\t%s\t%d\t\t%s\t%s\n", i+1, st, sr, o.Index, o.Path, tt)
	}
	fmt.Println("")
	if len(tasks) < 1 {
		return
	}
	c := -2
	for {
		fmt.Print("请输入序号进行状态取反(输入0返回):")
		for {
			fmt.Scanf("%d\n", &c)
			if c == 0 {
				return
			}
			if c-1 < 0 || c > len(tasks) {
				fmt.Print("参数错误,请重试(输入0返回):")
			} else {
				updateTaskStatus(c - 1)
				goto l1
			}
		}
	}
}

func add(x bool) {
	var url, path, y string
l1:
	fmt.Print("请输入地址(输入r返回):")
	fmt.Scanf("%s\n", &url)
	if url == "r" {
		return
	}
	var oa *api.LuzhiboAPI
	for {
		url = strings.ToLower(url)
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		oa = api.New(url)
		if oa == nil {
			fmt.Print("不支持的地址,请重新输入(输入r返回):")
			fmt.Scanf("%s\n", &url)
			if url == "r" {
				return
			}
		} else {
			break
		}
	}
	id, _, err := oa.GetRoomInfo()
	if err != nil {
		fmt.Print("获取地址信息失败")
		goto l1
	}
	tp := fmt.Sprintf("[%s]%s_%s", oa.Site, id, time.Now().Format("20060102150405"))
	fmt.Printf("请输入保存路径(输入r返回)[%s]:", tp)
	fmt.Scanf("%s\n", &path)
	if path == "r" {
		return
	}
	if path == "" {
		path = tp
	}
	if addTask(oa, path, x, false) {
		fmt.Print("添加成功,启动?(Y/n):")
		fmt.Scanf("%s\n", &y)
		if y == "" || y == "y" {
			startTask(len(tasks) - 1)
			fmt.Println("启动成功")
		}
	} else {
		fmt.Print("添加失败,重试?(Y/n):")
		fmt.Scanf("%s\n", &y)
		if y == "" || y == "y" {
			goto l1
		}
	}
}

func setConsoleTitle() {
	fmt.Printf("\033]0;%s\007", title)
}

func openWebUI(hta bool) {
	u := "http://localhost"
	if hta {

	}
	if port != 80 {
		u = fmt.Sprintf("%s:%d", u, port)
	}
	browser.OpenURL(u)
}

func quit() {
	for _, v := range tasks {
		v.worker.Stop()
	}
	os.Exit(0)
}
