package command

import (
	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/pkg/msend"
)

// Done updates a task to be marked done
type Done struct {
	doner *process.Update
}

func (d *Done) Cmd() string { return "done" }

func NewDone(localId int, conf *configuration.Configuration) (*Done, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Done{}, err
	}

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))
	fields := process.UpdateFields{"done": "true"}
	tId, err := findId(localId, local)
	if err != nil {
		return &Done{}, err
	}

	updater := process.NewUpdate(local, disp, tId, fields)

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
