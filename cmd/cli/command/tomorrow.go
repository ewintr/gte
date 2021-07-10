package command

import (
	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

// Tomorrow lists all tasks that are due tomorrow
type Tomorrow struct {
	tomorrower *process.List
}

func NewTomorrow(conf *configuration.Configuration) (*Tomorrow, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Tomorrow{}, err
	}

	reqs := process.ListReqs{
		Due: task.Today.Add(1),
	}
	tomorrower := process.NewList(local, reqs)

	return &Tomorrow{
		tomorrower: tomorrower,
	}, nil
}

func (t *Tomorrow) Do() string {
	res, err := t.tomorrower.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if len(res.Tasks) == 0 {
		return "nothing to do tomorrow\n"
	}

	return format.FormatTaskTable(res.Tasks)
}