package server

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"github.com/ehazlett/interlock/server/utils"
)

func (s *Server) getBackends(services []swarm.Service, tasks []swarm.Task) ([]*configurationapi.Backend, error) {
	svcs := map[string]swarm.Service{}
	for _, service := range services {
		svcs[service.ID] = service
	}

	backends := []*configurationapi.Backend{}
	for _, task := range tasks {
		service, ok := svcs[task.ServiceID]
		if !ok {
			return nil, fmt.Errorf("unable to find service spec for task service: %s", task.ServiceID)
		}
		b := &configurationapi.Backend{
			Name:  service.Spec.Name,
			Hosts: utils.Hosts(service.Spec),
		}

		backends = append(backends, b)
	}

	return backends, nil
}
