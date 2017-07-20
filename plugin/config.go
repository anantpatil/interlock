package plugin

import (
	"encoding/json"
	"os"

	"github.com/ehazlett/interlock/api/services/configuration"
	"github.com/ehazlett/interlock/api/types"
	"golang.org/x/net/context"
)

const (
	configPath = "/etc/interlock/config"
)

// LoadServiceConfig loads the service configuration from the well known location
func LoadServiceConfig() (*types.ServiceConfig, error) {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &types.ServiceConfig{}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (p *Plugin) Config(ctx context.Context) (*configuration.ConfigResponse, error) {
	config, err := p.client.Config(ctx, &configuration.ConfigRequest{
		ServiceCluster: p.serviceConfig.ServiceCluster,
	})
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (p *Plugin) Endpoint() string {
	if p.serviceConfig == nil {
		return ""
	}

	return p.serviceConfig.Endpoint
}
