package config

import (
	typesapi "github.com/ehazlett/interlock/api/types"
)

// Config is the top level configuration
type Config struct {
	ListenAddr       string
	GRPCAddr         string
	DockerURL        string
	TLSCACert        string
	TLSCert          string
	TLSKey           string
	AllowInsecure    bool
	EnableMetrics    bool
	PollInterval     string
	ProxyImage       string
	ProxyImageArgs   []string
	EndpointOverride string

	PluginConfig *typesapi.PluginConfig
}
