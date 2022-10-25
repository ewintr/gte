package storage

import (
	"time"

	"ewintr.nl/gte/internal/task"
	"github.com/google/uuid"
)

// Memory is an in memory implementation of LocalRepository
type Memory struct {
	tasks          map[string]*task.LocalTask
	latestFetch    time.Time
	latestDispatch time.Time
}

func NewMemory(initTasks ...*task.Task) *Memory {
	tasks := map[string]*task.LocalTask{}
	id := 1
	for _, t := range initTasks {
		tasks[t.Id] = &task.LocalTask{
			Task:        *t,
			LocalUpdate: &task.LocalUpdate{},
			LocalId:     id,
		}
		id++
	}

	return &Memory{
		tasks: tasks,
	}
}

func (m *Memory) LatestSyncs() (time.Time, time.Time, error) {
	return m.latestFetch, m.latestDispatch, nil
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
	m.latestFetch = time.Now()

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

func (m *Memory) SetLocalUpdate(id string, update *task.LocalUpdate) error {
	m.tasks[id].LocalStatus = task.STATUS_UPDATED
	m.tasks[id].LocalUpdate = update

	return nil
}

func (m *Memory) MarkDispatched(localId int) error {
	t, _ := m.FindByLocalId(localId)
	m.tasks[t.Id].LocalStatus = task.STATUS_DISPATCHED
	m.latestDispatch = time.Now()

	return nil
}

func (m *Memory) Add(update *task.LocalUpdate) (*task.LocalTask, error) {
	var used []int
	for _, t := range m.tasks {
		used = append(used, t.LocalId)
	}

	tsk := &task.LocalTask{
		Task: task.Task{
			Id:      uuid.New().String(),
			Version: 0,
			Folder:  task.FOLDER_NEW,
		},
		LocalId:     NextLocalId(used),
		LocalStatus: task.STATUS_UPDATED,
		LocalUpdate: update,
	}

	m.tasks[tsk.Id] = tsk

	return tsk, nil
}
