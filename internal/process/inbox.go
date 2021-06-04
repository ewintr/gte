package process

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrInboxProcess = errors.New("could not process inbox")

	inboxLock sync.Mutex
)

type Inbox struct {
	taskRepo *task.TaskRepo
}

type InboxResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewInbox(repo *task.TaskRepo) *Inbox {
	return &Inbox{
		taskRepo: repo,
	}
}

func (inbox *Inbox) Process() (*InboxResult, error) {
	inboxLock.Lock()
	defer inboxLock.Unlock()

	start := time.Now()

	tasks, err := inbox.taskRepo.FindAll(task.FOLDER_INBOX)
	if err != nil {
		return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
	}

	var cleanupNeeded bool
	for _, t := range tasks {
		if err := inbox.taskRepo.Update(t); err != nil {
			return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
		}
		cleanupNeeded = true
	}
	if cleanupNeeded {
		if err := inbox.taskRepo.CleanUp(); err != nil {
			return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
		}
	}

	return &InboxResult{
		Duration: time.Since(start).String(),
		Count:    len(tasks),
	}, nil
}
