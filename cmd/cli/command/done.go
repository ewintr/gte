package command

import (
	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
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

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))
	localTask, err := local.FindByLocalId(localId)
	if err != nil {
		return &Done{}, err
	}

	update := &task.LocalUpdate{
		ForVersion: localTask.Version,
		Fields:     []string{task.FIELD_DONE},
		Done:       true,
	}
	updater := process.NewUpdate(local, disp, localTask.Id, update)

	return &Done{
		doner: updater,
	}, nil
}

func (d *Done) Do() string {
	err := d.doner.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return "message sent\n"
}
