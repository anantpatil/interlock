package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	typesapi "github.com/ehazlett/interlock/api/types"
)

// ParseConfig returns a Config object from a raw string config TOML
func ParseConfig(data string) (*Config, error) {
	var cfg Config
	if _, err := toml.Decode(data, &cfg); err != nil {
		return nil, err
	}

	for _, p := range cfg.Plugins {
		if p.Config == nil {
			p.Config = &typesapi.PluginConfig{}
		}

		// check image
		if p.ProxyImage == "" {
			return nil, fmt.Errorf("ProxyImage must be specified")
		}

		// setup defaults for missing config entries
		SetConfigDefaults(p.Config)
	}

	return &cfg, nil
}

// SetConfigDefaults sets the defaults for the plugin config
func SetConfigDefaults(c *typesapi.PluginConfig) {
	if c.Version == "" {
		c.Version = "1"
	}

	if c.User == "" {
		c.User = "www-data"
	}

	if c.MaxConnections == 0 {
		c.MaxConnections = 1024
	}

	if c.Port == 0 {
		c.Port = 80
	}

	if c.WorkerProcesses == 0 {
		c.WorkerProcesses = 1
	}

	if c.RlimitNoFile == 0 {
		c.RlimitNoFile = 65535
	}

	if c.SendTimeout == 0 {
		c.SendTimeout = 600
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = 600
	}

	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 600
	}

	if c.ServerTimeout == 0 {
		c.ServerTimeout = 600
	}

	if c.ClientTimeout == 0 {
		c.ClientTimeout = 600
	}

	if c.SslCiphers == "" {
		c.SslCiphers = "HIGH:!aNULL:!MD5"
	}

	if c.SslProtocols == "" {
		c.SslProtocols = "SSLv3 TLSv1 TLSv1.1 TLSv1.2"
	}

	if c.AdminUser == "" {
		c.AdminUser = "admin"
	}

	if c.AdminPass == "" {
		c.AdminPass = ""
	}

	if c.SslDefaultDhParam == 0 {
		c.SslDefaultDhParam = 1024
	}

	if c.SslVerify == "" {
		c.SslVerify = "required"
	}
}
