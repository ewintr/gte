package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrRecurProcess = errors.New("could not generate tasks from recurrer")
)

type Recur struct {
	taskRepo       *task.TaskRepo
	taskDispatcher *task.Dispatcher
	daysAhead      int
}

type RecurResult struct {
	Count int
}

func NewRecur(repo *task.TaskRepo, disp *task.Dispatcher, daysAhead int) *Recur {
	return &Recur{
		taskRepo:       repo,
		taskDispatcher: disp,
		daysAhead:      daysAhead,
	}
}

func (recur *Recur) Process() (*RecurResult, error) {
	tasks, err := recur.taskRepo.FindAll(task.FOLDER_RECURRING)
	if err != nil {
		return &RecurResult{}, fmt.Errorf("%w: %v", ErrRecurProcess, err)
	}

	rDate := task.Today.AddDays(recur.daysAhead)
	var count int
	for _, t := range tasks {
		if t.RecursOn(rDate) {
			newTask, err := t.GenerateFromRecurrer(rDate)
			if err != nil {
				return &RecurResult{}, fmt.Errorf("%w: %v", ErrRecurProcess, err)
			}
			if err := recur.taskDispatcher.Dispatch(newTask); err != nil {
				return &RecurResult{}, fmt.Errorf("%w: %v", ErrRecurProcess, err)
			}
			count++
		}
	}

	return &RecurResult{
		Count: count,
	}, nil
}
