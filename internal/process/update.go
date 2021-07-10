package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/storage"
)

var (
	ErrUpdateTask = errors.New("could not update task")
)

// Update dispatches an updated version of a task
type Update struct {
	local   storage.LocalRepository
	disp    *storage.Dispatcher
	taskId  string
	updates UpdateFields
}

type UpdateFields map[string]string

func NewUpdate(local storage.LocalRepository, disp *storage.Dispatcher, taskId string, updates UpdateFields) *Update {
	return &Update{
		local:   local,
		disp:    disp,
		taskId:  taskId,
		updates: updates,
	}
}

func (u *Update) Process() error {
	task, err := u.local.FindById(u.taskId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	for k, v := range u.updates {
		switch k {
		case "done":
			if v == "true" {
				task.Done = true
			}
		}
	}

	if err := u.disp.Dispatch(task); err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	return nil
}
