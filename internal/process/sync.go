package process

import (
	"errors"
	"fmt"
	"time"

	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrSyncProcess = errors.New("could not sync local repository")
)

// Sync fetches all tasks in regular folders from the remote repository and overwrites what is stored locally
type Sync struct {
	remote *storage.RemoteRepository
	local  storage.LocalRepository
}

type SyncResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewSync(remote *storage.RemoteRepository, local storage.LocalRepository) *Sync {
	return &Sync{
		remote: remote,
		local:  local,
	}
}

func (s *Sync) Process() (*SyncResult, error) {
	start := time.Now()

	tasks := []*task.Task{}
	for _, folder := range task.KnownFolders {
		if folder == task.FOLDER_INBOX {
			continue
		}
		folderTasks, err := s.remote.FindAll(folder)
		if err != nil {
			return &SyncResult{}, fmt.Errorf("%w: %v", ErrSyncProcess, err)
		}

		for _, t := range folderTasks {
			tasks = append(tasks, t)
		}
	}

	if err := s.local.SetTasks(tasks); err != nil {
		return &SyncResult{}, fmt.Errorf("%w: %v", ErrSyncProcess, err)
	}

	return &SyncResult{
		Duration: time.Since(start).String(),
		Count:    len(tasks),
	}, nil
}
