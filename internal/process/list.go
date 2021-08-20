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
	Folder        string
	Project       string
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

	folders := []string{task.FOLDER_NEW, task.FOLDER_PLANNED, task.FOLDER_UNPLANNED}
	if l.reqs.Folder != "" {
		folders = []string{l.reqs.Folder}
	}

	var potentialTasks []*task.LocalTask
	for _, folder := range folders {
		folderTasks, err := l.local.FindAllInFolder(folder)
		if err != nil {
			return &ListResult{}, fmt.Errorf("%w: %v", ErrListProcess, err)
		}
		for _, ft := range folderTasks {
			potentialTasks = append(potentialTasks, ft)
		}
	}

	if l.reqs.Due.IsZero() && l.reqs.Project == "" {
		return &ListResult{
			Tasks: potentialTasks,
		}, nil
	}

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
