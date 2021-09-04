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

// Update updates a local task
type Update struct {
	local  storage.LocalRepository
	taskId string
	update *task.LocalUpdate
}

func NewUpdate(local storage.LocalRepository, taskId string, update *task.LocalUpdate) *Update {
	return &Update{
		local:  local,
		taskId: taskId,
		update: update,
	}
}

func (u *Update) Process() error {
	tsk, err := u.local.FindById(u.taskId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}
	tsk.AddUpdate(u.update)
	if err := u.local.SetLocalUpdate(tsk.Id, tsk.LocalUpdate); err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateTask, err)
	}

	return nil
}
