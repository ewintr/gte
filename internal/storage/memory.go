package storage

import (
	"time"

	"git.ewintr.nl/gte/internal/task"
)

// Memory is an in memory implementation of LocalRepository
type Memory struct {
	tasks      []*task.Task
	latestSync time.Time
	localIds   map[string]int
}

func NewMemory() *Memory {
	return &Memory{
		tasks:    []*task.Task{},
		localIds: map[string]int{},
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
		m.setLocalId(t.Id)
	}
	m.tasks = nTasks
	m.latestSync = time.Now()

	return nil
}

func (m *Memory) setLocalId(id string) {
	used := []int{}
	for _, id := range m.localIds {
		used = append(used, id)
	}

	next := NextLocalId(used)
	m.localIds[id] = next
}

func (m *Memory) FindAllInFolder(folder string) ([]*task.LocalTask, error) {
	tasks := []*task.LocalTask{}
	for _, t := range m.tasks {
		if t.Folder == folder {
			tasks = append(tasks, &task.LocalTask{
				Task:    *t,
				LocalId: m.localIds[t.Id],
			})
		}
	}

	return tasks, nil
}

func (m *Memory) FindAllInProject(project string) ([]*task.LocalTask, error) {
	tasks := []*task.LocalTask{}
	for _, t := range m.tasks {
		if t.Project == project {
			tasks = append(tasks, &task.LocalTask{
				Task:    *t,
				LocalId: m.localIds[t.Id],
			})
		}
	}

	return tasks, nil
}

func (m *Memory) FindById(id string) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if t.Id == id {
			return &task.LocalTask{
				Task:    *t,
				LocalId: m.localIds[t.Id],
			}, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}

func (m *Memory) FindByLocalId(localId int) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if m.localIds[t.Id] == localId {
			return &task.LocalTask{
				Task:    *t,
				LocalId: localId,
			}, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}
