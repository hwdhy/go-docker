package subsystem

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
	apply bool
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", cgroupPath, err)
		return err
	}
	if res.MemoryLimit != "" {
		m.apply = true
		//设置cgroup内存限制
		err = ioutil.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
		if err != nil {
			logrus.Errorf("set mem.limit_in_bytes err:%v", err)
			return err
		}
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsystemCgroupPath)
}

func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if m.apply {
		subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
		if err != nil {
			return err
		}
		taskPath := path.Join(subsystemCgroupPath, "tasks")
		err = ioutil.WriteFile(taskPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks, path: %s, pid: %d,err: %v", taskPath, pid, err)
			return err
		}
	}
	return nil
}
