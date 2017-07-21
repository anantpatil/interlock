package config

import (
	typesapi "github.com/ehazlett/interlock/api/types"
)

// Config is the top level configuration
type Config struct {
	ListenAddr    string
	GRPCAddr      string
	DockerURL     string
	TLSCACert     string
	TLSCert       string
	TLSKey        string
	AllowInsecure bool
	EnableMetrics bool
	PollInterval  string
	// GRPC endpoint override for plugins
	EndpointOverride string
	Plugins          []*Plugin
}

type Plugin struct {
	// Image to use for the plugin
	Image          string
	Args           []string
	ServiceCluster string
	Config         *typesapi.PluginConfig
}
