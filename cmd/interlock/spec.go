package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	typesapi "github.com/ehazlett/interlock/api/types"
	"github.com/ehazlett/interlock/config"
	"github.com/sirupsen/logrus"
)

var cmdSpec = cli.Command{
	Name:   "spec",
	Usage:  "generate a configuration file",
	Action: specAction,
}

func specAction(c *cli.Context) {
	pluginConfig := &typesapi.PluginConfig{
		ConfigPath: "/etc/proxy.conf",
	}

	config.SetConfigDefaults(pluginConfig)

	cfg := &config.Config{
		ListenAddr:    ":8080",
		GRPCAddr:      ":8081",
		DockerURL:     "unix:///var/run/docker.sock",
		EnableMetrics: true,
		PluginConfig:  pluginConfig,
		ProxyImage:    "ehazlett/interlock-plugin-nginx:latest",
	}

	if err := toml.NewEncoder(os.Stdout).Encode(cfg); err != nil {
		logrus.Fatal(err)
	}
}
