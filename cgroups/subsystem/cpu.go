package subsystem

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
	apply bool
}

func (c *CpuSubSystem) Name() string {
	return "cpu"
}

func (c *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", cgroupPath, err)
		return err
	}
	if res.CpuShare != "" {
		c.apply = true
		err = ioutil.WriteFile(path.Join(subsystemCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644)
		if err != nil {
			logrus.Errorf("failed to write file cpu.shares,err: %v", err)
			return err
		}
	}
	return nil
}

func (c *CpuSubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsystemCgroupPath)
}

func (c *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if c.apply {
		subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
		if err != nil {
			return err
		}

		taskPath := path.Join(subsystemCgroupPath, "tasks")
		err = ioutil.WriteFile(taskPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks, path: %s, pid: %d, err: %v", taskPath, pid, err)
			return err
		}
	}
	return nil
}
