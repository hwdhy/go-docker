package container

import (
	"docker_demo/common"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"syscall"
)

// 调用初始化函数，创建一个隔离namespace进程的Command，
func NewParentProcess(tty bool, volume string, containerName string, imageName string, envs []string) (*exec.Cmd, *os.File) {
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
	} else {
		//日志写入文件
		logDir := path.Join(common.DefaultContainerInfoPath, containerName)
		if _, err := os.Stat(logDir); err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				logrus.Errorf("mkdir container log, err:%v", err)
			}
		}
		logFileName := path.Join(logDir, common.ContainerLogFileName)
		file, err := os.Create(logFileName)
		if err != nil {
			logrus.Errorf("create log file err:%v", err)
		}
		cmd.Stdout = file
	}

	cmd.ExtraFiles = []*os.File{
		readPipe,
	}
	// 创建工作空间
	cmd.Env = append(cmd.Env, envs...)
	err := NewWorkSpace(volume, containerName, imageName)
	if err != nil {
		logrus.Errorf("new work space err: %v", err)
	}
	cmd.Dir = common.Merge + containerName

	return cmd, writePipe
}
