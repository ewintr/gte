package process

import (
	"errors"
	"fmt"

	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
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

func (s *Send) Process() (int, error) {
	tasks, err := s.local.FindAll()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrSendTasks, err)
	}

	var count int
	for _, t := range tasks {
		if t.LocalStatus != task.STATUS_UPDATED {
			continue
		}

		t.ApplyUpdate()
		if err := s.disp.Dispatch(&t.Task); err != nil {
			return 0, fmt.Errorf("%w: %v", ErrSendTasks, err)
		}
		if err := s.local.MarkDispatched(t.LocalId); err != nil {
			return 0, fmt.Errorf("%w: %v", ErrSendTasks, err)
		}

		count++
	}

	return count, nil
}
