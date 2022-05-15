package cgroups

import (
	"docker_demo/cgroups/subsystem"
	"github.com/sirupsen/logrus"
)

// CGroupManager 资源限制管理器
type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{
		path,
	}
}

func (c *CGroupManager) Set(res *subsystem.ResourceConfig) {
	for _, subSystem := range subsystem.Subsystems {
		err := subSystem.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("set %s err：%v", subSystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, subSystem := range subsystem.Subsystems {
		err := subSystem.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply task, err:%v", err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, subSystem := range subsystem.Subsystems {
		err := subSystem.Remove(c.Path)
		if err != nil {
			logrus.Errorf("remove %s err: %v", subSystem.Name(), err)
		}
	}
}
