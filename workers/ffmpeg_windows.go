//+build windows

package workers

import (
	nurl "net/url"
	"os/exec"
	"syscall"
)

func NewFFmpeg(url, path string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-y", "-i", url, "-vcodec", "copy", "-acodec", "copy", path)
	if Proxy != "" {
		_, err := nurl.Parse(url)
		if err == nil {
			cmd = exec.Command("ffmpeg", "-y", "-i", url, "-vcodec", "copy", "-acodec", "copy", "-http_proxy", Proxy, path)
		}
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
