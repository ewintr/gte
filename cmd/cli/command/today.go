package command

import (
	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

// Today lists all task that are due today or past their due date
type Today struct {
	local   storage.LocalRepository
	todayer *process.List
}

func NewToday(conf *configuration.Configuration) (*Today, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Today{}, err
	}
	reqs := process.ListReqs{
		Due:           task.Today,
		IncludeBefore: true,
	}
	todayer := process.NewList(local, reqs)

	return &Today{
		local:   local,
		todayer: todayer,
	}, nil
}

func (t *Today) Do() string {
	res, err := t.todayer.Process()
	if err != nil {
		return format.FormatError(err)
	}
	if len(res.Tasks) == 0 {
		return "nothing left\n"
	}

	return format.FormatTaskTable(t.local, res.Tasks)
}
