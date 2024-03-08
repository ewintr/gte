package command

import (
	"sort"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
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
		Due:          task.Today().Add(1),
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
	sort.Sort(task.ByDefault(res.Tasks))
	cols := []format.Column{format.COL_ID, format.COL_STATUS, format.COL_ACTION, format.COL_PROJECT}

	return format.FormatTaskTable(res.Tasks, cols)
}
