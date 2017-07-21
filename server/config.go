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
