package container

import (
	"docker_demo/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

// 创建容器运行时目录
func NewWorkSpace(volume string, containerName string, imageName string) error {
	// 1. 创建只读层
	err := createReadOnlyLayer(imageName)
	if err != nil {
		logrus.Errorf("create read only layer err: %v", err)
		return err
	}
	// 2. 创建读写层
	err = createUpperLayer(containerName)
	if err != nil {
		logrus.Errorf("create write layer err: %v", err)
		return err
	}
	// 3. 创建工作层
	err = createWorkLayer(containerName)
	if err != nil {
		logrus.Errorf("create work layer err: %v", err)
		return err
	}
	// 4. 创建挂载点
	err = CreateMountPoint(containerName, imageName)
	if err != nil {
		logrus.Errorf("create mount point err: %v", err)
		return err
	}
	// 5. 设置宿主机与容器文件映射
	mountVolume(containerName, imageName, volume)
	return nil
}

// 根据镜像创建只读层
func createReadOnlyLayer(imageName string) error {
	imagePath := path.Join(common.RootPath, imageName)
	_, err := os.Stat(imagePath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(imagePath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir image path err: %v", err)
			return err
		}
	}

	imageTarPath := path.Join(common.RootPath, fmt.Sprintf("%s.tar", imageName))
	if _, err := exec.Command("tar", "-xvf", imageTarPath, "-C", imagePath).CombinedOutput(); err != nil {
		logrus.Errorf("tar busybox.tar, err:%v", err)
		return err
	}
	return nil
}

// 创建读写层
func createUpperLayer(containerName string) error {
	upperLayerPath := path.Join(common.RootPath, common.Upper, containerName)
	_, err := os.Stat(upperLayerPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(upperLayerPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir upper layer, err: %v", err)
			return err
		}
	}
	return nil
}

// 创建工作层
func createWorkLayer(containerName string) error {
	workLayerPath := path.Join(common.RootPath, common.Work, containerName)
	_, err := os.Stat(workLayerPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(workLayerPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir work layer, err: %v", err)
			return err
		}
	}
	return nil
}

// CreateMountPoint 创建挂载点
func CreateMountPoint(containerName string, imageName string) error {
	mergePath := path.Join(common.Merge, containerName)
	_, err := os.Stat(mergePath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(mergePath, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir merge path err: %v", err)
			return err
		}
	}

	lowerDir := common.RootPath + imageName
	upperDir := common.RootPath + common.Upper + "/" + containerName
	workDir := common.RootPath + common.Work + "/" + containerName
	dirs := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	logrus.Info("mount", " -t", " overlay", " -o ", dirs, " overlay ", mergePath)
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "overlay", mergePath)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("merge cmd run err: %v", err)
		return err
	}
	return nil
}

// 宿主机和容器文件映射
func mountVolume(containerName, imageName, volume string) {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			// 创建宿主机中文件路径
			parentPath := volumes[0]
			if _, err := os.Stat(parentPath); err != nil && os.IsNotExist(err) {
				if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
					logrus.Errorf("mkdir parent path err: %v", err)
				}
			}

			// 创建容器内挂载点
			containerPath := volumes[1]
			containerVolumePath := path.Join(common.Merge, containerName, containerPath)
			if _, err := os.Stat(containerVolumePath); err != nil && os.IsNotExist(err) {
				if err := os.MkdirAll(containerVolumePath, os.ModePerm); err != nil {
					logrus.Errorf("mkdir volume path path: %s, err: %v", containerVolumePath, err)
				}
			}

			cmd := exec.Command("mount", "--bind", parentPath, containerVolumePath)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				logrus.Errorf("mount cmd run, err:%v", err)
			}
		}
	}
}

// 删除容器工作空间
func DeleteWorkSpace(containerName, volume string) error {
	// 1。 卸载挂载点
	mergePath := path.Join(common.Merge, containerName)
	err := unMountPoint(mergePath)
	if err != nil {
		return err
	}
	// 2. 删除读写层
	err = deleteWriteLayer(path.Join(common.RootPath, common.Upper, containerName))
	if err != nil {
		return err
	}
	// 3. 删除工作层
	err = deleteWorkLayer(path.Join(common.RootPath, common.Work, containerName))
	if err != nil {
		return err
	}

	// 4. 删除宿主机和文件系统映射
	deleteVolume(containerName, volume)
	return nil
}

// 卸载挂载点
func unMountPoint(mergePath string) error {
	if _, err := exec.Command("umount", mergePath).CombinedOutput(); err != nil {
		logrus.Errorf("umount merge path(%s) err: %v", mergePath, err)
		return err
	}

	err := os.RemoveAll(mergePath)
	if err != nil {
		logrus.Errorf("remove merge path err: %v", err)
		return err
	}
	return nil
}

// 删除读写层
func deleteWriteLayer(upperPath string) error {
	writerLayerPath := path.Join(upperPath, common.Upper)
	return os.RemoveAll(writerLayerPath)
}

// 删除工作层
func deleteWorkLayer(workPath string) error {
	workLayerPath := path.Join(workPath, common.Work)
	return os.RemoveAll(workLayerPath)
}

// 删除宿主机和文件系统映射
func deleteVolume(containerName, volume string) {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			containerPath := path.Join(common.Merge, containerName, volumes[1])
			if _, err := exec.Command("umount", containerPath).CombinedOutput(); err != nil {
				logrus.Errorf("umount container path(%s) err: %v", containerPath, err)
			}
		}
	}
}
