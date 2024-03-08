package command

import (
	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
)

type Update struct {
	updater *process.Update
}

func NewUpdate(localId int, conf *configuration.Configuration, cmdArgs []string) (*Update, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Update{}, err
	}

	update, err := ParseTaskFieldArgs(cmdArgs)
	if err != nil {
		return &Update{}, err
	}
	localTask, err := local.FindByLocalId(localId)
	if err != nil {
		return &Update{}, err
	}
	update.ForVersion = localTask.Version

	updater := process.NewUpdate(local, localTask.Id, update)

	return &Update{
		updater: updater,
	}, nil
}

func (u *Update) Do() string {
	if err := u.updater.Process(); err != nil {
		return format.FormatError(err)
	}

	return ""
}
