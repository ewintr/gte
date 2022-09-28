package command

import (
	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
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
