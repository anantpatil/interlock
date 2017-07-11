package nginx

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"github.com/ehazlett/interlock/plugin"
)

func init() {
	logrus.Debug("registering nginx")
	plugin.Register(&plugin.Registration{
		ID: "nginx",
		Init: func(c *plugin.InitContext) (interface{}, error) {
			return NewPlugin(c)
		},
	})
}

type Nginx struct {
}

func NewPlugin(c *plugin.InitContext) (*Nginx, error) {
	return &Nginx{}, nil
}

func (p *Nginx) Configure(services []swarm.Service) error {
	return nil
}

func (p *Nginx) Reload() error {
	return nil
}
