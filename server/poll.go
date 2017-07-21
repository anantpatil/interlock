package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	configurationapi "github.com/ehazlett/interlock/api/services/configuration"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (s *Server) poll() error {
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	optFilters := filters.NewArgs()
	optFilters.Add("desired-state", "running")
	opts := types.TaskListOptions{
		Filters: optFilters,
	}
	tasks, err := client.TaskList(context.Background(), opts)
	if err != nil {
		return errors.Wrap(err, "poll: unable to get tasks")
	}

	taskIDs := []string{}
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	sort.Strings(taskIDs)

	data, err := json.Marshal(taskIDs)
	if err != nil {
		return errors.Wrap(err, "unable to marshal task IDs")
	}

	h := sha256.New()
	h.Write(data)
	sum := hex.EncodeToString(h.Sum(nil))

	if sum != s.contentHash {
		// trigger update
		logrus.WithFields(logrus.Fields{
			"hash": sum,
		}).Info("update detected")
		s.contentHash = sum

		// TODO: build backend config and send to client
		s.currentConfig = &configurationapi.Config{
			Version: sum,
		}
	}

	return nil
}
