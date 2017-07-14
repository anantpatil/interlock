package server

import (
	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"golang.org/x/net/context"
)

func (s *Server) Config(ctx context.Context, req *configurationapi.ConfigRequest) (*configurationapi.ConfigResponse, error) {
	return &configurationapi.ConfigResponse{
		Config: s.currentConfig,
	}, nil
}
