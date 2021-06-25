package storage

import (
	"time"

	"git.ewintr.nl/gte/internal/task"
)

type LocalRepository interface {
	LatestSync() (time.Time, error)
	SetTasks(tasks []*task.Task) error
	FindAllInFolder(folder string) ([]*task.Task, error)
	FindAllInProject(project string) ([]*task.Task, error)
}
