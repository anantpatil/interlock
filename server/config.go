package server

import (
	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"golang.org/x/net/context"
)

func (s *Server) PluginConfig(ctx context.Context, req *configurationapi.PluginConfigRequest) (*configurationapi.PluginConfigResponse, error) {
	c, ok := s.plugins[req.ServiceCluster]
	if !ok {
		return nil, ErrServiceClusterConfigDoesNotExist
	}
	return &configurationapi.PluginConfigResponse{
		Config:       s.currentConfig,
		PluginConfig: c.Config,
	}, nil
}

func (s *Server) UpdateProxyConfig(ctx context.Context, req *configurationapi.UpdateProxyConfigRequest) (*configurationapi.UpdateProxyConfigResponse, error) {
	// TODO: create new proxy config
	_, version, err := s.createProxyServiceConfig(req.ServiceCluster, req.Data)
	if err != nil {
		return nil, err
	}

	// TODO: update proxy service
	if err := s.configureProxyService(req.ServiceCluster, version); err != nil {
		return nil, err
	}

	return &configurationapi.UpdateProxyConfigResponse{}, nil
}
