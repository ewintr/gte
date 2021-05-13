package task

import (
	"errors"
	"fmt"
	"strconv"

	"git.ewintr.nl/gte/pkg/mstore"
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
	if err := tr.Add(t); err != nil {
		return err
	}

	// remove old
	if err := tr.mstore.Remove(t.Message); err != nil {
		return fmt.Errorf("%w: %s", ErrMStoreError, err)
	}

	t.Current = false

	return nil
}

func (tr *TaskRepo) Add(t *Task) error {
	if t == nil {
		return ErrInvalidTask
	}

	if err := tr.mstore.Add(t.Folder, t.FormatSubject(), t.FormatBody()); err != nil {
		return fmt.Errorf("%w: %v", ErrMStoreError, err)
	}

	return nil
}

// Cleanup removes older versions of tasks
func (tr *TaskRepo) CleanUp() error {
	// loop through folders, get all task version info
	type msgInfo struct {
		Version int
		Message *mstore.Message
	}
	msgsSet := make(map[string][]msgInfo)

	for _, folder := range knownFolders {
		msgs, err := tr.mstore.Messages(folder)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrMStoreError, err)
		}

		for _, msg := range msgs {
			id, _ := FieldFromBody(FIELD_ID, msg.Body)
			versionStr, _ := FieldFromBody(FIELD_VERSION, msg.Body)
			version, _ := strconv.Atoi(versionStr)
			if _, ok := msgsSet[id]; !ok {
				msgsSet[id] = []msgInfo{}
			}
			msgsSet[id] = append(msgsSet[id], msgInfo{
				Version: version,
				Message: msg,
			})
		}
	}

	// determine which ones need to be gone
	var tobeRemoved []*mstore.Message
	for _, mInfos := range msgsSet {
		maxVersion := 0
		for _, mInfo := range mInfos {
			if mInfo.Version > maxVersion {
				maxVersion = mInfo.Version
			}
		}
		for _, mInfo := range mInfos {
			if mInfo.Version < maxVersion {
				tobeRemoved = append(tobeRemoved, mInfo.Message)
			}
		}
	}

	// remove them
	for _, msg := range tobeRemoved {
		if err := tr.mstore.Remove(msg); err != nil {
			return err
		}
	}

	return nil
}
