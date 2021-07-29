package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrUpdateTask = errors.New("could not update tsk")
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
	tsk, err := u.local.FindById(u.taskId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	for k, v := range u.updates {
		switch k {
		case task.FIELD_DONE:
			if v == "true" {
				tsk.Done = true
			}
		case task.FIELD_DUE:
			tsk.Due = task.NewDateFromString(v)
		case task.FIELD_ACTION:
			tsk.Action = v
		case task.FIELD_PROJECT:
			tsk.Project = v
		}
	}

	if err := u.disp.Dispatch(tsk); err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	return nil
}
