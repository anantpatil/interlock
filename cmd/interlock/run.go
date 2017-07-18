package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ehazlett/interlock/config"
	"github.com/ehazlett/interlock/server"
	"github.com/ehazlett/interlock/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultConfig = `ListenAddr = ":8080"
GRPCAddr = ":8081"
DockerURL = "unix:///var/run/docker.sock"
EnableMetrics = true
ProxyImage = "ehazlett/interlock-plugin-nginx:latest"
`
)

var cmdRun = cli.Command{
	Name:   "run",
	Usage:  "run interlock",
	Action: runAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config file",
			Value: "",
		},
	},
}

func runAction(c *cli.Context) error {
	logrus.Infof("interlock %s", version.FullVersion())

	var data string
	if envCfg := os.Getenv("INTERLOCK_CONFIG"); envCfg != "" {
		logrus.Debug("loading config from environment")
		data = envCfg
	}

	if configPath := c.String("config"); configPath != "" && data == "" {
		logrus.Debugf("loading config: file=%s", configPath)

		d, err := ioutil.ReadFile(configPath)
		switch {
		case os.IsNotExist(err):
			logrus.Errorf("Missing Interlock configuration: file=%s", configPath)
			logrus.Error("Use the run --config option to set a custom location for the configuration file")
			logrus.Error("Examples of an Interlock configuration file: url=https://github.com/ehazlett/interlock/tree/master/docs/examples")
			return fmt.Errorf("config not found: file=%s", configPath)
		case err == nil:
			data = string(d)
		default:
			return err
		}
	}

	if data == "" {
		logrus.Error("Examples of Interlock configuration: url=https://github.com/ehazlett/interlock/blob/master/docs/configuration.md")
		return fmt.Errorf("You must specify a config from a file or environment variable")
	}

	config, err := config.ParseConfig(data)
	if err != nil {
		return err
	}

	srv, err := server.NewServer(config)
	if err != nil {
		return err
	}

	if err := srv.Run(); err != nil {
		return err
	}

	return nil
}
