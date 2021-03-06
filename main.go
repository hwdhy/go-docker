package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"

	_ "docker_demo/nsenter"
)

const usage = "go-docker"

func main() {
	app := cli.NewApp()
	app.Name = "go-docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		runCommand,
		initCommand,
		logCommand,
		listCommand,
		commitCommand,
		execCommand,
		stopCommand,
		removeCommand,
		networkCommand,
	}

	app.Before = func(c *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
