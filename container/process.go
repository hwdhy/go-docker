package container

import (
	"docker_demo/common"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 调用初始化函数，创建一个隔离namespace进程的Command，
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, _ := os.Pipe()
	//调用自身，传入init参数，执行initCommand
	cmd := exec.Command("/proc/self/exe", "init")
	//设置隔离信息
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	//交互模式，日志写入控制台
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{
		readPipe,
	}
	// 创建工作空间
	err := NewWorkSpace(common.RootPath, common.Merge, volume)
	if err != nil {
		logrus.Errorf("new work space err: %v", err)
	}
	cmd.Dir = common.Merge

	return cmd, writePipe
}
