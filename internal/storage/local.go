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
}
