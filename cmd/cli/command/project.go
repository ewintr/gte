package command

import (
	"sort"
	"strings"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
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

	sort.Sort(task.ByDefault(res.Tasks))
	cols := []format.Column{format.COL_ID, format.COL_STATUS, format.COL_DUE, format.COL_ACTION}

	return format.FormatTaskTable(res.Tasks, cols)
}
