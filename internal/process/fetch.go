package process

import (
	"errors"
	"fmt"
	"time"

	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
)

var (
	ErrFetchProcess = errors.New("could not fetch tasks")
)

// Fetch fetches all tasks in regular folders from the remote repository and overwrites what is stored locally
type Fetch struct {
	remote *storage.RemoteRepository
	local  storage.LocalRepository
}

type FetchResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewFetch(remote *storage.RemoteRepository, local storage.LocalRepository) *Fetch {
	return &Fetch{
		remote: remote,
		local:  local,
	}
}

func (s *Fetch) Process() (*FetchResult, error) {
	start := time.Now()

	tasks := []*task.Task{}
	for _, folder := range task.KnownFolders {
		if folder == task.FOLDER_INBOX {
			continue
		}
		folderTasks, err := s.remote.FindAll(folder)
		if err != nil {
			return &FetchResult{}, fmt.Errorf("%w: %v", ErrFetchProcess, err)
		}

		for _, t := range folderTasks {
			tasks = append(tasks, t)
		}
	}

	if err := s.local.SetTasks(tasks); err != nil {
		return &FetchResult{}, fmt.Errorf("%w: %v", ErrFetchProcess, err)
	}

	return &FetchResult{
		Duration: time.Since(start).String(),
		Count:    len(tasks),
	}, nil
}
