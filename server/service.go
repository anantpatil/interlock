package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/ehazlett/interlock/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Server) getProxyService() (*swarm.Service, error) {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	optFilters := filters.NewArgs()
	optFilters.Add("label", proxyServiceLabel)
	opts := types.ServiceListOptions{
		Filters: optFilters,
	}
	services, err := client.ServiceList(context.Background(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get services")
	}

	if len(services) == 0 {
		return nil, nil
	}

	if len(services) > 1 {
		return nil, fmt.Errorf("found more than one proxy service: %+v", services)
	}

	return &services[0], nil
}

func (s *Server) proxyServiceConfigExists() (bool, error) {
	svc, err := s.getProxyService()
	if err != nil {
		return false, err
	}

	return svc != nil, nil
}

func (s *Server) removeProxyServiceConfig() error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	optFilters := filters.NewArgs()
	optFilters.Add("name", proxyServiceConfigName)
	opts := types.ConfigListOptions{
		Filters: optFilters,
	}

	cfgs, err := client.ConfigList(context.Background(), opts)
	if err != nil {
		return err
	}

	if len(cfgs) > 1 {
		return fmt.Errorf("more than one config found: %+v", cfgs)
	}

	cfg := cfgs[0]
	if err := client.ConfigRemove(context.Background(), cfg.ID); err != nil {
		return err
	}

	return nil
}

func (s *Server) createProxyServiceConfig() (swarm.Config, error) {
	cfg := swarm.Config{}

	client, err := getDockerClient(s.cfg)
	if err != nil {
		return cfg, err
	}
	defer client.Close()

	serviceConfigData, err := json.Marshal(s.serviceConfig)
	if err != nil {
		return cfg, err
	}

	spec := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name: proxyServiceConfigName,
			Labels: map[string]string{
				"version": version.FullVersion(),
			},
		},
		Data: serviceConfigData,
	}

	id := ""

	proxyServiceExists, err := s.proxyServiceConfigExists()
	if err != nil {
		return cfg, err
	}

	if !proxyServiceExists {
		resp, err := client.ConfigCreate(context.Background(), spec)
		if err != nil {
			return cfg, err
		}
		id = resp.ID
	}

	c, _, err := client.ConfigInspectWithRaw(context.Background(), id)
	if err != nil {
		return cfg, err
	}

	return c, nil
}

func (s *Server) createProxyService() error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	serviceConfig, err := s.createProxyServiceConfig()
	if err != nil {
		return err
	}

	taskSpec := swarm.TaskSpec{
		ContainerSpec: &swarm.ContainerSpec{
			Image: s.cfg.ProxyImage,
			//Secrets: []*swarm.SecretReference{
			//    {
			//        File: &swarm.SecretReferenceFileTarget{
			//    	Name: "/etc/interlock/cert.pem",
			//        },
			//    },
			//},
			Configs: []*swarm.ConfigReference{
				{
					File: &swarm.ConfigReferenceFileTarget{
						Name: "/etc/interlock/config",
						UID:  "0",
						GID:  "0",
						Mode: 0644,
					},
					ConfigName: proxyServiceConfigName,
					ConfigID:   serviceConfig.ID,
				},
			},
		},
		RestartPolicy: &swarm.RestartPolicy{
			Condition: swarm.RestartPolicyConditionAny,
		},
	}
	if len(s.cfg.ProxyImageArgs) > 0 {
		taskSpec.ContainerSpec.Args = s.cfg.ProxyImageArgs
	}
	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Labels: map[string]string{
				proxyServiceLabel: "",
				"version":         version.FullVersion(),
			},
		},
		TaskTemplate: taskSpec,
	}

	svc, err := client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error creating service: %+v", spec)
	}

	logrus.WithFields(logrus.Fields{
		"id": svc.ID,
	}).Debug("proxy service")

	return nil
}
