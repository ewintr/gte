package command

import (
	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
)

// Today lists all task that are due today or past their due date
type Today struct {
	lister *process.List
}

func NewToday(conf *configuration.Configuration) (*Today, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Today{}, err
	}
	reqs := process.ListReqs{
		Due:           task.Today,
		IncludeBefore: true,
		ApplyUpdates:  true,
	}
	lister := process.NewList(local, reqs)

	return &Today{
		lister: lister,
	}, nil
}

func (t *Today) Do() string {
	res, err := t.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return format.FormatTaskTable(res.Tasks)
}
