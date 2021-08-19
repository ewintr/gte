package command

import (
	"errors"
	"strings"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

var (
	ErrUnknownFolder = errors.New("unknown folder")
)

// List lists all the tasks in a project or a folder
type List struct {
	local  storage.LocalRepository
	lister *process.List
}

func (l *List) Cmd() string { return "list" }

func NewList(conf *configuration.Configuration, cmdArgs []string) (*List, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &List{}, err
	}
	if len(cmdArgs) < 2 {
		return &List{}, ErrInvalidAmountOfArgs
	}

	reqs, err := parseReqs(cmdArgs[0], cmdArgs[1])
	if err != nil {
		return &List{}, err
	}
	lister := process.NewList(local, reqs)

	return &List{
		local:  local,
		lister: lister,
	}, nil
}

func (l *List) Do() string {
	res, err := l.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if len(res.Tasks) == 0 {
		return "no tasks there\n"
	}

	return format.FormatTaskTable(l.local, res.Tasks)
}

func parseReqs(kind, item string) (process.ListReqs, error) {
	item = strings.ToLower(item)
	switch kind {
	case "folder":
		for _, folder := range task.KnownFolders {
			if item == strings.ToLower(folder) {
				return process.ListReqs{
					Folder: folder,
				}, nil
			}
		}
		return process.ListReqs{}, ErrUnknownFolder
	case "project":
		return process.ListReqs{
			Project: item,
		}, nil
	}

	return process.ListReqs{}, process.ErrInvalidReqs
}
