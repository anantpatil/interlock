package plugin

import (
	"github.com/ehazlett/interlock/api/services/configuration"
	"golang.org/x/net/context"
)

func (p *Plugin) Config(ctx context.Context) (*configuration.Config, error) {
	config, err := p.client.Config(ctx, &configuration.ConfigRequest{})
	if err != nil {
		return nil, err
	}

	return config.Config, nil
}
