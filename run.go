package main

import (
	"docker_demo/cgroups"
	"docker_demo/cgroups/subsystem"
	"docker_demo/common"
	"docker_demo/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(cmdArray []string, tty bool, res *subsystem.ResourceConfig, volume string, containerName string) {
	//创建隔离namespace的cmd
	parentProcess, writePipe := container.NewParentProcess(tty, volume, containerName)
	if parentProcess == nil {
		logrus.Errorf("failed to new parent process")
		return
	}

	if err := parentProcess.Start(); err != nil {
		logrus.Errorf("parent start failed, err: %v", err)
		return
	}
	//添加资源限制
	cGroupManager := cgroups.NewCGroupManager("go-docker")
	//进程退出 清除资源限制文件
	defer cGroupManager.Destroy()
	//设置资源限制
	cGroupManager.Set(res)
	//将容器进程加入到各个subsystem挂载对应的cgroup中
	cGroupManager.Apply(parentProcess.Process.Pid)

	//设置初始化命令
	setInitCommand(cmdArray, writePipe)
	//等待父进程结束
	err := parentProcess.Wait()
	if err != nil {
		logrus.Errorf("parent wait err : %v", err)
	}

	err = container.DeleteWorkSpace(common.RootPath, common.Merge, volume)
	if err != nil {
		logrus.Errorf("delete work pace err: %v", err)
	}
}

// 设置初始化cmd
func setInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	logrus.Infof("command all is %s", command)

	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
