package command

import (
	"errors"
	"fmt"
	"strings"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

var (
	ErrFieldAlreadyUsed = errors.New("field was already used")
)

type Update struct {
	updater *process.Update
}

func (u *Update) Cmd() string { return "update" }

func NewUpdate(localId int, conf *configuration.Configuration, cmdArgs []string) (*Update, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Update{}, err
	}

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))
	fields, err := ParseTaskFieldArgs(cmdArgs)
	if err != nil {
		return &Update{}, err
	}
	tId, err := findId(localId, local)
	if err != nil {
		return &Update{}, err
	}

	updater := process.NewUpdate(local, disp, tId, fields)

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

func ParseTaskFieldArgs(args []string) (process.UpdateFields, error) {
	result := process.UpdateFields{}

	var action []string
	for _, f := range args {
		split := strings.SplitN(f, ":", 2)
		if len(split) == 2 {
			switch split[0] {
			case "project":
				if _, ok := result[task.FIELD_PROJECT]; ok {
					return process.UpdateFields{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_PROJECT)
				}
				result[task.FIELD_PROJECT] = split[1]
			case "due":
				if _, ok := result[task.FIELD_DUE]; ok {
					return process.UpdateFields{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_DUE)
				}
				result[task.FIELD_DUE] = split[1]
			}
		} else {
			action = append(action, f)
		}
	}

	if len(action) > 0 {
		result[task.FIELD_ACTION] = strings.Join(action, " ")
	}

	return result, nil
}
