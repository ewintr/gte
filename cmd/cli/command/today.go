package command

import (
	"sort"

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
		Due:           task.Today(),
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

	sort.Sort(task.ByDefault(res.Tasks))
	cols := []format.Column{format.COL_ID, format.COL_STATUS, format.COL_ACTION, format.COL_PROJECT}

	return format.FormatTaskTable(res.Tasks, cols)
}
