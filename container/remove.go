package container

import (
	"docker_demo/common"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

//删除容器
func RemoveContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info err: %v", err)
		return
	}

	if info.Status != common.Stop {
		logrus.Errorf("can't remove runnign container.")
		return
	}

	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	err = os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("remove container dir: %s, err: %v", dir, err)
	}
}
