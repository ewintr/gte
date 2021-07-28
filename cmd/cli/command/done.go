package command

import (
	"fmt"

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

func NewDone(id int, conf *configuration.Configuration) (*Done, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Done{}, err
	}

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))
	fields := process.UpdateFields{"done": "true"}
	localIds, err := local.LocalIds()
	if err != nil {
		return &Done{}, err
	}
	var tId string
	for remoteId, localId := range localIds {
		if localId == id {
			tId = remoteId
			break
		}
	}
	if tId == "" {
		return &Done{}, fmt.Errorf("could not find task")
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
