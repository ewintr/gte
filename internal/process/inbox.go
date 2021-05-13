package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrInboxProcess = errors.New("could not process inbox")
)

type Inbox struct {
	taskRepo *task.TaskRepo
}

type InboxResult struct {
	Count int
}

func NewInbox(repo *task.TaskRepo) *Inbox {
	return &Inbox{
		taskRepo: repo,
	}
}

func (inbox *Inbox) Process() (*InboxResult, error) {
	tasks, err := inbox.taskRepo.FindAll(task.FOLDER_INBOX)
	if err != nil {
		return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
	}

	var cleanupNeeded bool
	for _, t := range tasks {
		if t.Dirty {
			if err := inbox.taskRepo.Update(t); err != nil {
				return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
			}
			cleanupNeeded = true
		}
	}
	if cleanupNeeded {
		if err := inbox.taskRepo.CleanUp(); err != nil {
			return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
		}
	}

	return &InboxResult{
		Count: len(tasks),
	}, nil
}
