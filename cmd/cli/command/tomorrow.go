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
	lister *process.List
}

func NewTomorrow(conf *configuration.Configuration) (*Tomorrow, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Tomorrow{}, err
	}

	reqs := process.ListReqs{
		Due:          task.Today.Add(1),
		ApplyUpdates: true,
	}
	lister := process.NewList(local, reqs)

	return &Tomorrow{
		lister: lister,
	}, nil
}

func (t *Tomorrow) Do() string {
	res, err := t.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return format.FormatTaskTable(res.Tasks)
}
