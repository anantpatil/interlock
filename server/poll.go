package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
)

func (s *Server) poll() error {
	logrus.Debug("poller tick")
	client, err := getDockerClient(s.cfg)
	if err != nil {
		return err
	}

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
		}).Debug("triggering update")
		s.contentHash = sum
	}

	return nil
}
