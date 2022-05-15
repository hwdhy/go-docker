package subsystem

var (
	Subsystems = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)

// ResourceConfig 资源限制配置
type ResourceConfig struct {
	//内存限制
	MemoryLimit string
	//CPU时间片权重
	CpuShare string
	//CPU核数
	CpuSet string
}

type Subsystem interface {
	Name() string
	Set(cgroupPath string, res *ResourceConfig) error
	Remove(cgroupPath string) error
	Apply(cgroupPath string, pid int) error
}
