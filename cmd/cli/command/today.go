package command

import (
	"fmt"

	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

type Today struct {
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
		todayer: todayer,
	}, nil
}

func (t *Today) Do() string {
	res, err := t.todayer.Process()
	if err != nil {
		return FormatError(err)
	}
	if len(res.Tasks) == 0 {
		return "nothing left\n"
	}

	var msg string
	for _, t := range res.Tasks {
		msg += fmt.Sprintf("%s - %s\n", t.Project, t.Action)
	}

	return msg
}
