package server

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"github.com/pkg/errors"
)

func (s *Server) updateConfiguration() error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	services, err := client.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return errors.Wrap(err, "update: unable to get services")
	}

	optFilters := filters.NewArgs()
	optFilters.Add("desired-state", "running")
	opts := types.TaskListOptions{
		Filters: optFilters,
	}
	tasks, err := client.TaskList(context.Background(), opts)
	if err != nil {
		return errors.Wrap(err, "update: unable to get tasks")
	}

	// TODO: for each service ensure proxy service is connected to specified network for access
	//proxyNetworks := []string{}

	// build backend config and send to client
	backends, err := s.getBackends(services, tasks)
	if err != nil {
		return errors.Wrap(err, "update: unable to get backends")
	}

	data, err := json.Marshal(backends)
	if err != nil {
		return err
	}

	version := generateHash(data)

	s.currentConfig = &configurationapi.Config{
		Version:  version,
		Backends: backends,
	}

	return nil
}
