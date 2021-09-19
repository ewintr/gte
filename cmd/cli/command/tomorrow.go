package command

import (
	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
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
