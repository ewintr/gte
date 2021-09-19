package command

import (
	"strings"

	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
)

type Project struct {
	lister *process.List
}

func NewProject(conf *configuration.Configuration, cmdArgs []string) (*Project, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Project{}, err
	}
	if len(cmdArgs) < 1 {
		return &Project{}, ErrInvalidAmountOfArgs
	}
	reqs := process.ListReqs{
		Project:      strings.ToLower(cmdArgs[0]),
		ApplyUpdates: true,
	}
	lister := process.NewList(local, reqs)

	return &Project{
		lister: lister,
	}, nil
}

func (p *Project) Do() string {
	res, err := p.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if len(res.Tasks) == 0 {
		return "no tasks here\n\n"
	}

	return format.FormatTaskTable(res.Tasks)
}
