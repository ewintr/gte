package storage

import (
	"time"

	"git.ewintr.nl/gte/internal/task"
)

type localData struct {
	LocalId     int
	LocalUpdate *task.LocalUpdate
}

// Memory is an in memory implementation of LocalRepository
type Memory struct {
	tasks      []*task.Task
	latestSync time.Time
	localData  map[string]localData
}

func NewMemory() *Memory {
	return &Memory{
		tasks:     []*task.Task{},
		localData: map[string]localData{},
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
	for _, ld := range m.localData {
		used = append(used, ld.LocalId)
	}

	next := NextLocalId(used)
	m.localData[id] = localData{
		LocalId: next,
	}
}

func (m *Memory) FindAllInFolder(folder string) ([]*task.LocalTask, error) {
	tasks := []*task.LocalTask{}
	for _, t := range m.tasks {
		if t.Folder == folder {
			tasks = append(tasks, &task.LocalTask{
				Task:        *t,
				LocalId:     m.localData[t.Id].LocalId,
				LocalUpdate: m.localData[t.Id].LocalUpdate,
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
				Task:        *t,
				LocalId:     m.localData[t.Id].LocalId,
				LocalUpdate: m.localData[t.Id].LocalUpdate,
			})
		}
	}

	return tasks, nil
}

func (m *Memory) FindById(id string) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if t.Id == id {
			return &task.LocalTask{
				Task:        *t,
				LocalId:     m.localData[t.Id].LocalId,
				LocalUpdate: m.localData[t.Id].LocalUpdate,
			}, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}

func (m *Memory) FindByLocalId(localId int) (*task.LocalTask, error) {
	for _, t := range m.tasks {
		if m.localData[t.Id].LocalId == localId {
			return &task.LocalTask{
				Task:        *t,
				LocalId:     localId,
				LocalUpdate: m.localData[t.Id].LocalUpdate,
			}, nil
		}
	}

	return &task.LocalTask{}, ErrTaskNotFound
}

func (m *Memory) SetLocalUpdate(localId int, localUpdate *task.LocalUpdate) error {
	t, err := m.FindByLocalId(localId)
	if err != nil {
		return err
	}

	m.localData[t.Id] = localData{
		LocalId:     localId,
		LocalUpdate: localUpdate,
	}

	return nil
}
