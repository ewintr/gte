package command

import (
	"go-mod.ewintr.nl/gte/cmd/cli/format"
	"go-mod.ewintr.nl/gte/internal/configuration"
	"go-mod.ewintr.nl/gte/internal/process"
	"go-mod.ewintr.nl/gte/internal/storage"
)

// New sends an action to the NEW folder so it can be updated to a real task later
type New struct {
	newer *process.New
}

func NewNew(conf *configuration.Configuration, cmdArgs []string) (*New, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &New{}, err
	}

	update, err := ParseTaskFieldArgs(cmdArgs)
	if err != nil {
		return &New{}, err
	}

	return &New{
		newer: process.NewNew(local, update),
	}, nil
}

func (n *New) Do() string {
	if err := n.newer.Process(); err != nil {
		return format.FormatError(err)
	}

	return ""
}
