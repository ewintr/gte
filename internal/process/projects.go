package process

import (
	"errors"
	"fmt"
	"sort"

	"git.ewintr.nl/gte/internal/storage"
)

var (
	ErrCouldNotFetchProjects = errors.New("could not fetch projects")
)

type Projects struct {
	local storage.LocalRepository
}

func NewProjects(local storage.LocalRepository) *Projects {
	return &Projects{
		local: local,
	}
}

func (p *Projects) Process() ([]string, error) {
	allTasks, err := p.local.FindAll()
	if err != nil {
		return []string{}, fmt.Errorf("%w: %v", ErrCouldNotFetchProjects, err)
	}

	knownMap := map[string]bool{}
	for _, t := range allTasks {
		knownMap[t.Project] = true
	}
	known := []string{}
	for p := range knownMap {
		known = append(known, p)
	}
	sort.Strings(known)

	return known, nil
}
