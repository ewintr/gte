package storage

import (
	"errors"
	"time"

	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrTaskNotFound = errors.New("task was not found")
)

type LocalRepository interface {
	LatestSync() (time.Time, error)
	SetTasks(tasks []*task.Task) error
	FindAllInFolder(folder string) ([]*task.Task, error)
	FindAllInProject(project string) ([]*task.Task, error)
	FindById(id string) (*task.Task, error)
	FindByLocalId(id int) (*task.Task, error)
	LocalIds() (map[string]int, error)
}

func NextLocalId(used []int) int {
	if len(used) == 0 {
		return 1
	}

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
