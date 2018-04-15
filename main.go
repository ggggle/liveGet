//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	//"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Baozisoftware/GoldenDaemon"
	"github.com/ggggle/liveGet/api"
	"github.com/ggggle/liveGet/api/getters"
	"github.com/ggggle/liveGet/workers"
)

const ver = 2018041500
const p = "录直播"

var port = 22216

var nhta *bool

var htaproc *os.Process

var nt *bool

var proxy *string
var logFile, _ = os.Create("lzb.log")
//var logBuf = bytes.NewBufferString("")
var logger = log.New(logFile, "", log.LstdFlags)

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	flag.Bool("d", false, "启用后台运行(仅非Windows有效)")
	nt = flag.Bool("nt", false, "启用无终端交互模式(仅非Windows有效)")
	proxy = flag.String("proxy", "", "代理服务器(如:\"http://127.0.0.1:8888\".)")
	pieceSize := flag.Int64("size", workers.ONE_PIECE_SIZE, "分片段大小，单位MB")
	flag.Parse()
	port = *p
	s := ":" + strconv.Itoa(port)
	if runtime.GOOS != "windows" {
		GoldenDaemon.RegisterTrigger("d", "-nt")
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(time.Second * 5)
		d, f := filepath.Split(os.Args[0])
		tp := filepath.Join(d, "."+f+".old")
		os.Remove(tp)
	}()
	getters.Proxy = *proxy
	workers.Proxy = *proxy
	api.Logger = logger
	workers.ONE_PIECE_SIZE = *pieceSize * 1024 * 1024
	logger.Printf("片段大小[%d]MB", *pieceSize)
	logger.Print("软件启动成功.")
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	if !*nt || runtime.GOOS == "windows" {
		time.Sleep(time.Second * 2)
		go startServer(s)
		if !*nopen {
			openWebUI(!*nhta)
		}
		cmd()
	} else {
		startServer(s)
	}
}
