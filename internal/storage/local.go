package storage

import (
	"errors"
	"sort"
	"time"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrTaskNotFound = errors.New("task was not found")
)

type LocalRepository interface {
	LatestSync() (time.Time, error)
	SetTasks(tasks []*task.Task) error
	FindAllInFolder(folder string) ([]*task.LocalTask, error)
	FindAllInProject(project string) ([]*task.LocalTask, error)
	FindById(id string) (*task.LocalTask, error)
	FindByLocalId(id int) (*task.LocalTask, error)
}

func NextLocalId(used []int) int {
	if len(used) == 0 {
		return 1
	}

	sort.Ints(used)
	usedMax := 1
	for _, u := range used {
		if u > usedMax {
			usedMax = u
		}
	}

	var limit int
	for limit = 1; limit <= len(used) || limit < usedMax; limit *= 10 {
	}

	newId := used[len(used)-1] + 1
	if newId < limit {
		return newId
	}

	usedMap := map[int]bool{}
	for _, u := range used {
		usedMap[u] = true
	}

	for i := 1; i < limit; i++ {
		if _, ok := usedMap[i]; !ok {
			return i
		}
	}

	return limit
}
