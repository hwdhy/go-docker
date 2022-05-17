package common

const (
	RootPath   = "/home/cater/"
	Merge      = "/home/cater/merge/"
	Lower      = "busybox"
	BusyBoxTar = "busybox.tar"
	Upper      = "upper"
	Work       = "work"
	BinPath    = "/bin/"
)

const (
	DefaultContainerInfoPath = "/var/run/docker/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"
)

const (
	Running = "running"
	Stop    = "stopped"
	Exit    = "exited"
)

const (
	EnvExecPid = "docker_pid"
	EnvExecCmd = "docker_cmd"
)

const (
	DefaultNetworkPath   = "/var/run/go-docker/network/network/"
	DefaultAllocatorPath = "/var/run/go-docker/network/ipam/subnet.json"
)
