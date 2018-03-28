//+build windows

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/dkua/go-ico"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/pkg/browser"
)

func cmd() {
	p := true
	resp, _ := http.Get(webuiHost() + "/favicon.ico")
	defer resp.Body.Close()
	imgs, _ := ico.DecodeAll(resp.Body)
	icon, _ := walk.NewIconFromImage(imgs.Image[6])
	ni, _ := walk.NewNotifyIcon()
	defer ni.Dispose()
	mw, _ := walk.NewMainWindow()
	ni.SetIcon(icon)
	ni.SetToolTip(fmt.Sprintf("正在\"%d\"处监听WebUI...", port))
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton && p {
			openWebUI(!*nhta)
		} else if button == walk.RightButton && p {
			p = false
			if walk.MsgBox(mw, "录直播", "退出后将停止所有正在运行的任务,确定退出?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) == win.IDYES {
				if htaproc != nil {
					htaproc.Kill()
				}
				mw.Close()
			}
			p = true
		}
	})
	ni.SetVisible(true)
	ni.ShowCustom("录直播 - 软件已启动...", "左键点击重新打开WebUI,右键点击退出本软件.")
	mw.Run()
}

func quit() {
	if htaproc != nil {
		htaproc.Kill()
	}
	for _, v := range tasks {
		v.worker.Stop()
	}
	os.Exit(0)
}

func openWebUI(hta bool) {
	if htaproc != nil {
		htaproc.Kill()
	}
	u := webuiHost()
	if runtime.GOOS == "windows" && hta && checkWin10() {
		cmd := exec.Command("mshta", u+"/hta")
		err := cmd.Start()
		if err != nil {
			browser.OpenURL(u)
		} else {
			htaproc = cmd.Process
		}
	} else {
		browser.OpenURL(u)
	}
}

func webuiHost() string {
	u := "http://localhost"
	if port != 80 {
		u = fmt.Sprintf("%s:%d", u, port)
	}
	return u
}

func checkWin10() (ret bool) {
	defer func() {
		if recover() != nil {
			ret = false
		}
	}()
	type OSVERSIONINFOW struct {
		dwOSVersionInfoSize uint32
		dwMajorVersion      uint32
		dwMinorVersion      uint32
		dwBuildNumber       uint32
		dwPlatformId        uint32
		szCSDVersion        [128]byte
	}
	var v OSVERSIONINFOW
	p := uintptr(unsafe.Pointer(&v))
	kernel32 := syscall.NewLazyDLL("ntdll.dll")
	c := kernel32.NewProc("RtlGetVersion")
	r, _, _ := c.Call(p)
	if r == 0 {
		ret = v.dwMajorVersion == 10
	} else {
		ret = false
	}
	return
}
