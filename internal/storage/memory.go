package storage

import (
	"time"

	"git.ewintr.nl/gte/internal/task"
)

// Memory is an in memory implementation of LocalRepository
type Memory struct {
	tasks      []*task.Task
	latestSync time.Time
}

func NewMemory() *Memory {
	return &Memory{
		tasks: []*task.Task{},
	}
}

func (m *Memory) LatestSync() (time.Time, error) {
	return m.latestSync, nil
}

func (m *Memory) SetTasks(tasks []*task.Task) error {
	nTasks := []*task.Task{}
	for _, t := range tasks {
		nt := *t
		nt.Message = nil
		nTasks = append(nTasks, &nt)
	}
	m.tasks = nTasks
	m.latestSync = time.Now()

	return nil
}

func (m *Memory) FindAllInFolder(folder string) ([]*task.Task, error) {
	tasks := []*task.Task{}
	for _, t := range m.tasks {
		if t.Folder == folder {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (m *Memory) FindAllInProject(project string) ([]*task.Task, error) {
	tasks := []*task.Task{}
	for _, t := range m.tasks {
		if t.Project == project {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (m *Memory) FindById(id string) (*task.Task, error) {
	for _, t := range m.tasks {
		if t.Id == id {
			return t, nil
		}

	}

	return &task.Task{}, ErrTaskNotFound
}
