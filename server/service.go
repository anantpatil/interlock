package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (s *Server) getProxyService() (*swarm.Service, error) {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	optFilters := filters.NewArgs()
	optFilters.Add("label", "type="+proxyServiceLabel)
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

func (s *Server) getProxyServiceConfig(version string) (*swarm.Config, error) {
	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Debugf("checking service config")
	cfgs, err := s.getServiceConfigs()
	if err != nil {
		return nil, err
	}

	for _, cfg := range cfgs {
		if cfg.Spec.Labels["version"] == version {
			return &cfg, nil
		}
	}

	return nil, nil
}

func (s *Server) createProxyServiceConfig() (swarm.Config, string, error) {
	cfg := swarm.Config{}

	client, err := getDockerClient(s.cfg)
	if err != nil {
		return cfg, "", err
	}
	defer client.Close()

	serviceConfigData, err := json.Marshal(s.serviceConfig)
	if err != nil {
		return cfg, "", err
	}

	version := generateHash(serviceConfigData)

	spec := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name: proxyServiceConfigName + "." + version,
			Labels: map[string]string{
				"type":    proxyServiceConfigName,
				"version": version,
			},
		},
		Data: serviceConfigData,
	}

	config, err := s.getProxyServiceConfig(version)
	if err != nil {
		return cfg, "", err
	}

	if config == nil {
		if _, err := client.ConfigCreate(context.Background(), spec); err != nil {
			return cfg, "", err
		}
		c, err := s.getProxyServiceConfig(version)
		if err != nil {
			return cfg, "", err
		}

		config = c
	}

	return *config, version, nil
}

func (s *Server) configureProxyService() error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	serviceConfig, version, err := s.createProxyServiceConfig()
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
					ConfigID:   serviceConfig.ID,
					ConfigName: serviceConfig.Spec.Name,
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
				"type": proxyServiceLabel,
			},
		},
		TaskTemplate: taskSpec,
	}

	serviceID := ""

	svc, err := s.getProxyService()
	if err != nil {
		return err
	}

	if svc == nil {
		service, err := client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
		if err != nil {
			return errors.Wrapf(err, "error creating service: %+v", spec)
		}

		serviceID = service.ID
	} else {
		opts := types.ServiceUpdateOptions{}
		// update service to remove current config
		clearSpec := svc.Spec
		clearSpec.TaskTemplate.ContainerSpec.Configs = []*swarm.ConfigReference{}
		if _, err := client.ServiceUpdate(context.Background(), svc.ID, svc.Version, clearSpec, opts); err != nil {
			return err
		}

		// TODO: wait for service to be updated using UpdateStatus
		time.Sleep(time.Second * 1)

		// get updated service with new version
		updatedService, err := s.getProxyService()
		if err != nil {
			return err
		}

		// update service with new config
		spec.Annotations.Name = svc.Spec.Name

		if _, err := client.ServiceUpdate(context.Background(), svc.ID, updatedService.Version, spec, opts); err != nil {
			return err
		}

		time.Sleep(time.Second * 1)

		serviceID = updatedService.ID
	}

	// remove old configs
	if err := s.cleanProxyServiceConfigs(version); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"id": serviceID,
	}).Debug("proxy service")

	return nil
}

func (s *Server) getServiceConfigs() ([]swarm.Config, error) {
	cfgs := []swarm.Config{}
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return cfgs, err
	}
	defer client.Close()

	optFilters := filters.NewArgs()
	//optFilters.Add("label", "type="+proxyServiceConfigName)
	opts := types.ConfigListOptions{
		Filters: optFilters,
	}

	configs, err := client.ConfigList(context.Background(), opts)
	if err != nil {
		return cfgs, err
	}

	return configs, nil
}

func (s *Server) cleanProxyServiceConfigs(currentHash string) error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	cfgs, err := s.getServiceConfigs()
	if err != nil {
		return err
	}

	for _, cfg := range cfgs {
		if cfg.Spec.Labels["version"] != currentHash {
			logrus.WithFields(logrus.Fields{
				"id":   cfg.ID,
				"name": cfg.Spec.Annotations.Name,
			}).Debug("removing old service config")
			if err := client.ConfigRemove(context.Background(), cfg.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
