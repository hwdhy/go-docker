package main

import (
	"docker_demo/cgroups/subsystem"
	"docker_demo/common"
	"docker_demo/container"
	"docker_demo/network"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroup limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "docker volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "docker env",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container args")
		}

		res := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		var cmdArray []string
		for _, arg := range context.Args().Tail() {
			cmdArray = append(cmdArray, arg)
		}

		tty := context.Bool("it")
		detach := context.Bool("d")

		if tty && detach {
			return fmt.Errorf("it and d paramter can not both provided")
		}

		volume := context.String("v")
		containerName := context.String("name")
		net := context.String("net")

		imageName := context.Args().Get(0)
		envs := context.StringSlice("e")
		ports := context.StringSlice("p")

		// 入口函数
		Run(cmdArray, tty, res, volume, containerName, imageName, net, envs, ports)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container, Do not call it outside",
	Action: func(context *cli.Context) error {
		logrus.Infof("init come on")
		return container.RunContainerInitProcess()
	},
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "look container log",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}

		containerName := context.Args().Get(0)
		container.LookContainerLog(containerName)
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		// 如果环境变量里面有 PID,那么则什么都不执行
		pid := os.Getenv(common.EnvExecPid)
		if pid != "" {
			logrus.Infof("pid callback pid %s, gid: %d", pid, os.Getgid())
			return nil
		}
		if len(context.Args()) < 2 {
			return fmt.Errorf("missing container name or command")
		}

		var cmdArray []string
		for _, arg := range context.Args().Tail() {
			cmdArray = append(cmdArray, arg)
		}

		containerName := context.Args().Get(0)
		container.ExecContainer(containerName, cmdArray)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "look container log",
	Action: func(context *cli.Context) error {
		container.ListContainerInfo()
		return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "docker commit a container into image",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Usage: "export image path",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		imageName := context.Args().Get(0)
		imagePath := context.String("c")
		return container.CommitContainer(imageName, imagePath)
	},
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing stop container name")
		}
		containerName := context.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "rm a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing remove container name")
		}
		containerName := context.Args().Get(0)
		container.RemoveContainer(containerName)
		return nil
	},
}

var networkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				err := network.Init()
				if err != nil {
					logrus.Errorf("network init failed, err: %v", err)
					return err
				}
				// 创建网络
				err = network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				err := network.Init()
				if err != nil {
					logrus.Errorf("network init failed, err: %v", err)
					return err
				}

				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}

				err := network.Init()
				if err != nil {
					logrus.Errorf("network init failed, err: %v", err)
					return err
				}
				err = network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}
