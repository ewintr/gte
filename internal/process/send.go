package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrSendTasks = errors.New("could not send tasks")
)

// Send sends local tasks that need to be dispatched
type Send struct {
	local storage.LocalRepository
	disp  *storage.Dispatcher
}

func NewSend(local storage.LocalRepository, disp *storage.Dispatcher) *Send {
	return &Send{
		local: local,
		disp:  disp,
	}
}

func (s *Send) Process() error {
	tasks, err := s.local.FindAll()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendTasks, err)
	}

	for _, t := range tasks {
		if t.LocalStatus != task.STATUS_UPDATED {
			continue
		}

		t.ApplyUpdate()
		if err := s.disp.Dispatch(&t.Task); err != nil {
			return fmt.Errorf("%w: %v", ErrSendTasks, err)
		}
		if err := s.local.MarkDispatched(t.LocalId); err != nil {
			return fmt.Errorf("%w: %v", ErrSendTasks, err)
		}
	}

	return nil
}
