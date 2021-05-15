package process

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrRecurProcess = errors.New("could not generate tasks from recurrer")

	recurLock sync.Mutex
)

type Recur struct {
	taskRepo       *task.TaskRepo
	taskDispatcher *task.Dispatcher
	daysAhead      int
}

type RecurResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewRecur(repo *task.TaskRepo, disp *task.Dispatcher, daysAhead int) *Recur {
	return &Recur{
		taskRepo:       repo,
		taskDispatcher: disp,
		daysAhead:      daysAhead,
	}
}

func (recur *Recur) Process() (*RecurResult, error) {
	recurLock.Lock()
	defer recurLock.Unlock()

	start := time.Now()

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
		Duration: time.Since(start).String(),
		Count:    count,
	}, nil
}
