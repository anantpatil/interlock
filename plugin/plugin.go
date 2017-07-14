package plugin

import (
	"github.com/ehazlett/interlock/api/services/configuration"
	"google.golang.org/grpc"
)

type Plugin struct {
	client configuration.ConfigurationClient
}

func NewPlugin(addr string) (*Plugin, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := configuration.NewConfigurationClient(conn)

	return &Plugin{
		client: client,
	}, nil
}
