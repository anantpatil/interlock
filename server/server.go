package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock/config"
	"github.com/ehazlett/interlock/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	defaultPollInterval = time.Second * 2
)

type Server struct {
	cfg         *config.Config
	plugins     []plugin.Plugin
	metrics     *Metrics
	contentHash string
}

func NewServer(cfg *config.Config) (*Server, error) {
	plugins := plugin.Plugins()
	logrus.Debugf("plugins: %+v", plugins)

	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) Run() error {
	logrus.Info("starting server")
	if s.cfg.EnableMetrics {
		// start prometheus listener
		http.Handle("/metrics", prometheus.Handler())
	}

	if s.cfg.PollInterval != "" {
		// run background poller
		d, err := time.ParseDuration(s.cfg.PollInterval)
		if err != nil {
			return err
		}

		if d < defaultPollInterval {
			log.Warnf("poll interval too quick; defaulting to %v", defaultPollInterval)
			s.cfg.PollInterval = "2s"
			d = defaultPollInterval
		}

		// start poller
		t := time.NewTicker(d)
		go func() {
			for range t.C {
				if err := s.poll(); err != nil {
					logrus.Error(err)
				}
			}
		}()

	}

	if err := http.ListenAndServe(s.cfg.ListenAddr, nil); err != nil {
		return err
	}

	return nil
}
