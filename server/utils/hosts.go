package utils

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/ehazlett/interlock"
)

func Hosts(spec swarm.ServiceSpec) []string {
	hosts := []string{}

	for l, v := range spec.Labels {
		// this is for labels like interlock.hosts.1=foo.local
		if strings.Index(l, interlock.InterlockHostsLabel) > -1 {
			hosts = append(hosts, v)
		}
	}

	return hosts
}
