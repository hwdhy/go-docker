package container

import (
	"docker_demo/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"path"
)

// CommitContainer 导出容器
func CommitContainer(imageName string, imagePath string) error {
	if imagePath == "" {
		imagePath = common.RootPath
	}

	imageTar := path.Join(imagePath, fmt.Sprintf("%s.tar", imageName))
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", common.Merge, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar container image, file name: %s, err %v", imageName, err)
		return err
	}
	return nil
}
