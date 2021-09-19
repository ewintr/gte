package command

import (
	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
)

// Done updates a task to be marked done
type Done struct {
	doner *process.Update
}

func NewDone(localId int, conf *configuration.Configuration) (*Done, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Done{}, err
	}
	localTask, err := local.FindByLocalId(localId)
	if err != nil {
		return &Done{}, err
	}

	update := &task.LocalUpdate{
		ForVersion: localTask.Version,
		Fields:     []string{task.FIELD_DONE},
		Done:       true,
	}
	updater := process.NewUpdate(local, localTask.Id, update)

	return &Done{
		doner: updater,
	}, nil
}

func (d *Done) Do() string {
	err := d.doner.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return ""
}
