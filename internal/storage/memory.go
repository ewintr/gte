package storage

import (
	"time"

	"git.ewintr.nl/gte/internal/task"
)

// Memory is an in memory implementation of LocalRepository
type Memory struct {
	tasks      map[string]*task.LocalTask
	latestSync time.Time
}

func NewMemory() *Memory {
	return &Memory{
		tasks: map[string]*task.LocalTask{},
	}
}

func (m *Memory) LatestSync() (time.Time, error) {
	return m.latestSync, nil
}

func (m *Memory) SetTasks(tasks []*task.Task) error {
	var oldTasks []*task.LocalTask
	for _, ot := range m.tasks {
		oldTasks = append(oldTasks, ot)
	}

	newTasks := MergeNewTaskSet(oldTasks, tasks)

	m.tasks = map[string]*task.LocalTask{}
	for _, nt := range newTasks {
		m.tasks[nt.Id] = nt
	}
	m.latestSync = time.Now()

	return nil
}

func (m *Memory) FindAll() ([]*task.LocalTask, error) {
	tasks := []*task.LocalTask{}
	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m *Memory) FindById(id string) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if t.Id == id {
			return t, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}

func (m *Memory) FindByLocalId(localId int) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if t.LocalId == localId {
			return t, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}

func (m *Memory) SetLocalUpdate(tsk *task.LocalTask) error {
	m.tasks[tsk.Id] = tsk

	return nil
}
