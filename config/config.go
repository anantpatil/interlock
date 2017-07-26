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
	// Image to use for the plugin service
	Image string
	Args  []string
	// Image for the proxy service
	ProxyImage string
	ProxyArgs  []string
	// Config path in the service for the proxy config
	ProxyConfigPath string
	ServiceCluster  string
	Config          *typesapi.PluginConfig
}
