package process

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
)

var (
	ErrInboxProcess = errors.New("could not process inbox")

	inboxLock sync.Mutex
)

// Inbox processes all messages in INBOX in a remote repository
type Inbox struct {
	taskRepo *storage.RemoteRepository
}

type InboxResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewInbox(repo *storage.RemoteRepository) *Inbox {
	return &Inbox{
		taskRepo: repo,
	}
}

func (inbox *Inbox) Process() (*InboxResult, error) {
	inboxLock.Lock()
	defer inboxLock.Unlock()

	start := time.Now()

	// find tasks to be processed
	tasks, err := inbox.taskRepo.FindAll(task.FOLDER_INBOX)
	if err != nil {
		return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
	}

	// deduplicate
	taskKeys := map[string]*task.Task{}
	for _, newT := range tasks {
		existingT, ok := taskKeys[newT.Id]
		switch {
		case !ok:
			taskKeys[newT.Id] = newT
		case newT.Version >= existingT.Version:
			taskKeys[newT.Id] = newT
		}
	}
	tasks = []*task.Task{}
	for _, t := range taskKeys {
		tasks = append(tasks, t)
	}

	// split them
	doneTasks, updateTasks := []*task.Task{}, []*task.Task{}
	for _, t := range tasks {
		if t.Done {
			doneTasks = append(doneTasks, t)
			continue
		}
		updateTasks = append(updateTasks, t)
	}

	// remove
	if err := inbox.taskRepo.Remove(doneTasks); err != nil {
		return &InboxResult{}, fmt.Errorf("%w: %v", ErrInboxProcess, err)
	}

	// update
	var cleanupNeeded bool
	for _, t := range updateTasks {
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
