package process

import (
	"errors"
	"fmt"

	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
)

var (
	ErrNewTask = errors.New("could not add new task")
)

type New struct {
	local  storage.LocalRepository
	update *task.LocalUpdate
}

func NewNew(local storage.LocalRepository, update *task.LocalUpdate) *New {
	return &New{
		local:  local,
		update: update,
	}
}

func (n *New) Process() error {
	if _, err := n.local.Add(n.update); err != nil {
		return fmt.Errorf("%w: %v", ErrNewTask, err)
	}

	return nil
}
