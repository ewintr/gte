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
	local  storage.LocalRepository
	disp   *storage.Dispatcher
	taskId string
	update task.LocalUpdate
}

func NewUpdate(local storage.LocalRepository, disp *storage.Dispatcher, taskId string, update task.LocalUpdate) *Update {
	return &Update{
		local:  local,
		disp:   disp,
		taskId: taskId,
		update: update,
	}
}

func (u *Update) Process() error {
	tsk, err := u.local.FindById(u.taskId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}
	tsk.Apply(u.update)

	if err := u.disp.Dispatch(&tsk.Task); err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	return nil
}
