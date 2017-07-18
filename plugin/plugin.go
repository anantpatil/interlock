package plugin

import (
	"github.com/ehazlett/interlock/api/services/configuration"
	"github.com/ehazlett/interlock/api/types"
	"google.golang.org/grpc"
)

type Plugin struct {
	client        configuration.ConfigurationClient
	serviceConfig *types.ServiceConfig
}

func NewPlugin(addr string) (*Plugin, error) {
	p := &Plugin{}
	if addr == "" {
		cfg, err := LoadServiceConfig()
		if err != nil {
			return nil, err
		}

		p.serviceConfig = cfg

		if cfg != nil {
			addr = cfg.Endpoint
		}
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	p.client = configuration.NewConfigurationClient(conn)

	return p, nil
}
