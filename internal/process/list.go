package process

import (
	"errors"
	"fmt"

	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
)

var (
	ErrInvalidReqs = errors.New("could not make sense of requirements")
	ErrListProcess = errors.New("could not fetch task list")
)

// ListReqs specifies the requirements in AND fashion for a list of tasks
type ListReqs struct {
	Due           task.Date
	IncludeBefore bool
	Folder        string
	Project       string
	ApplyUpdates  bool
}

func (lr ListReqs) Valid() bool {
	switch {
	case lr.Folder != "":
		return true
	case lr.Project != "":
		return true
	case !lr.Due.IsZero():
		return true
	}

	return false
}

// List finds all tasks that satisfy the given requirements
type List struct {
	local storage.LocalRepository
	reqs  ListReqs
}

type ListResult struct {
	Tasks []*task.LocalTask
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

	potentialTasks, err := l.local.FindAll()
	if err != nil {
		return &ListResult{}, fmt.Errorf("%w: %v", ErrListProcess, err)
	}

	// updates
	if l.reqs.ApplyUpdates {
		for i := range potentialTasks {
			potentialTasks[i].ApplyUpdate()
		}
		var undoneTasks []*task.LocalTask
		for _, pt := range potentialTasks {
			if !pt.Done {
				undoneTasks = append(undoneTasks, pt)
			}
		}
		potentialTasks = undoneTasks
	}

	// folder
	if l.reqs.Folder != "" {
		var folderTasks []*task.LocalTask
		for _, pt := range potentialTasks {
			if pt.Folder == l.reqs.Folder {
				folderTasks = append(folderTasks, pt)
			}
		}

		potentialTasks = folderTasks
	}

	if l.reqs.Due.IsZero() && l.reqs.Project == "" {
		return &ListResult{
			Tasks: potentialTasks,
		}, nil
	}

	// project
	if l.reqs.Project != "" {
		var projectTasks []*task.LocalTask
		for _, pt := range potentialTasks {
			if pt.Project == l.reqs.Project {
				projectTasks = append(projectTasks, pt)
			}
		}

		potentialTasks = projectTasks
	}

	if l.reqs.Due.IsZero() {
		return &ListResult{
			Tasks: potentialTasks,
		}, nil
	}

	dueTasks := []*task.LocalTask{}
	for _, t := range potentialTasks {
		switch {
		case t.Due.IsZero():
			// skip
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
