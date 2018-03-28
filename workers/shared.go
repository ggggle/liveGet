package workers

import (
	"os"
	"path"
	"github.com/ggggle/luzhibo/api/getters"
)

//Worker 工作接口
type Worker interface {
	Start()                                                           //开始
	Stop()                                                            //停止
	Restart() (Worker, error)                                         //重新开始
	GetTaskInfo(bool) (int64, bool, int64, string, *getters.LiveInfo) //取状态
}

//WorkCompletedCallBack 工作完成回调
type WorkCompletedCallBack func(int64)

func createFile(filepath string) (file *os.File, err error) {
	err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err == nil {
		return os.Create(filepath)
	}
	return
}

var Proxy=""