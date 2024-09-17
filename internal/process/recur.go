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
	ErrRecurProcess = errors.New("could not generate tasks from recurrer")

	recurLock sync.Mutex
)

// Recur generates new tasks from a recurring task for a given day
type Recur struct {
	taskRepo       *storage.RemoteRepository
	taskDispatcher *storage.Dispatcher
	daysAhead      int
}

type RecurResult struct {
	Duration string `json:"duration"`
	Count    int    `json:"count"`
}

func NewRecur(repo *storage.RemoteRepository, disp *storage.Dispatcher, daysAhead int) *Recur {
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

	rDate := task.Today().AddDays(recur.daysAhead)
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
