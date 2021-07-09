package process

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrInvalidReqs = errors.New("could not make sense of requirements")
	ErrListProcess = errors.New("could not fetch task list")
)

// ListReqs specifies the requirements in AND fashion for a list of tasks
type ListReqs struct {
	Due           task.Date
	IncludeBefore bool
}

func (lr ListReqs) Valid() bool {
	return !lr.Due.IsZero()
}

// List finds all tasks that satisfy the given requirements
type List struct {
	local storage.LocalRepository
	reqs  ListReqs
}

type ListResult struct {
	Tasks []*task.Task
}

func NewList(local storage.LocalRepository, reqs ListReqs) *List {
	return &List{
		local: local,
		reqs:  reqs,
	}
}

func (l *List) Process() (*ListResult, error) {
	if !l.reqs.Valid() {
		return &ListResult{}, ErrInvalidReqs
	}

	potentialTasks, err := l.local.FindAllInFolder(task.FOLDER_PLANNED)
	if err != nil {
		return &ListResult{}, fmt.Errorf("%w: %v", ErrListProcess, err)
	}

	dueTasks := []*task.Task{}
	for _, t := range potentialTasks {
		switch {
		case t.Due.Equal(l.reqs.Due):
			dueTasks = append(dueTasks, t)
		case l.reqs.IncludeBefore && l.reqs.Due.After(t.Due):
			dueTasks = append(dueTasks, t)
		}
	}

	return &ListResult{
		Tasks: dueTasks,
	}, nil
}
