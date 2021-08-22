package command

import (
	"fmt"
	"strings"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

type Update struct {
	updater *process.Update
}

func NewUpdate(localId int, conf *configuration.Configuration, cmdArgs []string) (*Update, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Update{}, err
	}

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))
	update, err := ParseTaskFieldArgs(cmdArgs)
	if err != nil {
		return &Update{}, err
	}
	localTask, err := local.FindByLocalId(localId)
	if err != nil {
		return &Update{}, err
	}
	update.ForVersion = localTask.Version

	updater := process.NewUpdate(local, disp, localTask.Id, update)

	return &Update{
		updater: updater,
	}, nil
}

func (u *Update) Do() string {
	if err := u.updater.Process(); err != nil {
		return format.FormatError(err)
	}

	return "message sent\n"
}

func ParseTaskFieldArgs(args []string) (*task.LocalUpdate, error) {
	lu := &task.LocalUpdate{}

	action, fields := []string{}, []string{}
	for _, f := range args {
		split := strings.SplitN(f, ":", 2)
		if len(split) == 2 {
			switch split[0] {
			case "project":
				if lu.Project != "" {
					return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_PROJECT)
				}
				lu.Project = split[1]
				fields = append(fields, task.FIELD_PROJECT)
			case "due":
				if !lu.Due.IsZero() {
					return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_DUE)
				}
				lu.Due = task.NewDateFromString(split[1])
				fields = append(fields, task.FIELD_DUE)
			}
		} else {
			if len(f) > 0 {
				action = append(action, f)
			}
		}
	}

	if len(action) > 0 {
		lu.Action = strings.Join(action, " ")
		fields = append(fields, task.FIELD_ACTION)
	}

	lu.Fields = fields

	return lu, nil
}
