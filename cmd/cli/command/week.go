package command

import (
	"sort"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
)

type Week struct {
	lister *process.List
}

func NewWeek(conf *configuration.Configuration) (*Week, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Week{}, err
	}

	reqs := process.ListReqs{
		Due:           task.Today().Add(7),
		IncludeBefore: true,
		ApplyUpdates:  true,
	}
	return &Week{
		lister: process.NewList(local, reqs),
	}, nil
}

func (w *Week) Do() string {
	res, err := w.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	sort.Sort(task.ByDefault(res.Tasks))

	return format.FormatTaskTable(res.Tasks, format.COL_ALL)
}
