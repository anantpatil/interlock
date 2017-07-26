package server

import (
	"encoding/json"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
		return errors.Wrap(err, "poll: unable to marshal task IDs")
	}

	version := generateHash(data)

	if version != s.contentHash {
		// trigger update
		logrus.WithFields(logrus.Fields{
			"version": version,
		}).Info("update detected")
		s.contentHash = version

		if err := s.updateConfiguration(); err != nil {
			return errors.Wrap(err, "poll: unable to update configuration")
		}
	}

	return nil
}
