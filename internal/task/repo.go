package task

import (
	"errors"
	"fmt"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

var (
	ErrMStoreError    = errors.New("mstore gave error response")
	ErrInvalidTask    = errors.New("invalid task")
	ErrInvalidMessage = errors.New("task contains invalid message")
)

type TaskRepo struct {
	mstore mstore.MStorer
}

func NewRepository(ms mstore.MStorer) *TaskRepo {
	return &TaskRepo{
		mstore: ms,
	}
}

func (tr *TaskRepo) FindAll(folder string) ([]*Task, error) {
	msgs, err := tr.mstore.Messages(folder)
	if err != nil {
		return []*Task{}, fmt.Errorf("%w: %v", ErrMStoreError, err)
	}

	tasks := []*Task{}
	for _, msg := range msgs {
		if msg.Valid() {
			tasks = append(tasks, New(msg))
		}
	}

	return tasks, nil
}

func (tr *TaskRepo) Update(t *Task) error {
	if t == nil {
		return ErrInvalidTask
	}
	if !t.Current {
		return ErrOutdatedTask
	}
	if !t.Dirty {
		return nil
	}

	// add new
	if err := tr.mstore.Add(t.Folder, t.FormatSubject(), t.FormatBody()); err != nil {
		return fmt.Errorf("%w: %s", ErrMStoreError, err)
	}

	// remove old
	if err := tr.mstore.Remove(t.Message); err != nil {
		return fmt.Errorf("%w: %s", ErrMStoreError, err)
	}

	t.Current = false

	return nil
}

// Cleanup removes older versions of tasks
func (tr *TaskRepo) CleanUp() error {
	// loop through folders, get all tasks
	taskSet := make(map[string][]*Task)

	for _, folder := range knownFolders {
		tasks, err := tr.FindAll(folder)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if _, ok := taskSet[t.Id]; !ok {
				taskSet[t.Id] = []*Task{}
			}
			taskSet[t.Id] = append(taskSet[t.Id], t)
		}
	}

	// determine which ones need to be gone
	var tobeRemoved []*Task
	for _, tasks := range taskSet {
		maxVersion := 0
		for _, t := range tasks {
			if t.Version > maxVersion {
				maxVersion = t.Version
			}
		}

		for _, t := range tasks {
			if t.Version < maxVersion {
				tobeRemoved = append(tobeRemoved, t)
			}
		}
	}

	// remove them
	for _, t := range tobeRemoved {
		if err := tr.mstore.Remove(t.Message); err != nil {
			return err
		}
	}

	return nil
}
