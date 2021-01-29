package task

import (
	"errors"
	"fmt"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

var (
	ErrMStoreError = errors.New("mstore gave error response")
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
		return []*Task{}, err
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
