package plugin

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
)

var (
	ErrNoPluginID = errors.New("plugin: no id")
)

type PluginType string

type Registration struct {
	Type   PluginType
	ID     string
	Config interface{}
	Init   func(*InitContext) (interface{}, error)

	added bool
}

type Plugin interface {
	ID() string
	Configure(services []swarm.Service) error
	Reload() error
}

var register = struct {
	sync.Mutex
	r []*Registration
}{}

func Register(r *Registration) {
	logrus.WithFields(logrus.Fields{
		"id": r.ID,
	}).Debug("registering plugin")
	register.Lock()
	defer register.Unlock()
	if r.ID == "" {
		panic(ErrNoPluginID)
	}
	register.r = append(register.r, r)
}

func Plugins() (plugins []*Registration) {
	return register.r
}
