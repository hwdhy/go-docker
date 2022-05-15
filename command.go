package main

import (
	"docker_demo/cgroups/subsystem"
	"docker_demo/container"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container args")
		}
		tty := context.Bool("it")
		volume := context.String("v")
		containerName := context.String("name")

		res := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		// 入口函数
		Run(cmdArray, tty, res, volume, containerName)
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
