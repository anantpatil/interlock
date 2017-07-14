package server

import (
	"net"
	"net/http"
	"time"

	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"github.com/ehazlett/interlock/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultPollInterval = time.Second * 2
)

type Server struct {
	cfg           *config.Config
	metrics       *Metrics
	contentHash   string
	currentConfig *configurationapi.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) Run() error {
	logrus.Info("starting server")
	if s.cfg.EnableMetrics {
		// start prometheus listener
		http.Handle("/metrics", prometheus.Handler())
		go func() {
			if err := http.ListenAndServe(s.cfg.ListenAddr, nil); err != nil {
				logrus.Error("unable to start metric listener: %s", err)
			}
		}()

	}

	if s.cfg.PollInterval != "" {
		// run background poller
		d, err := time.ParseDuration(s.cfg.PollInterval)
		if err != nil {
			return err
		}

		if d < defaultPollInterval {
			log.Warnf("poll interval too quick; defaulting to %s", defaultPollInterval)
			s.cfg.PollInterval = "2s"
			d = defaultPollInterval
		}

		// start poller
		t := time.NewTicker(d)
		go func() {
			for range t.C {
				logrus.Debug("poll")
				if err := s.poll(); err != nil {
					logrus.Error(err)
				}
			}
		}()
	}

	grpcServer := grpc.NewServer()
	configurationapi.RegisterConfigurationServer(grpcServer, s)

	l, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		return err
	}

	logrus.Debug("starting GRPC server")
	grpcServer.Serve(l)

	return nil
}
