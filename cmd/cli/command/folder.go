package command

import (
	"strings"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

type Folder struct {
	lister *process.List
}

func NewFolder(conf *configuration.Configuration, cmdArgs []string) (*Folder, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Folder{}, err
	}
	if len(cmdArgs) < 1 {
		return &Folder{}, ErrInvalidAmountOfArgs
	}
	knownFolders := []string{task.FOLDER_NEW, task.FOLDER_PLANNED, task.FOLDER_UNPLANNED}
	var folder string
	for _, f := range knownFolders {
		if strings.ToLower(f) == strings.ToLower(cmdArgs[0]) {
			folder = f
			break
		}
	}
	if folder == "" {
		return &Folder{}, ErrUnknownFolder
	}

	reqs := process.ListReqs{
		Folder:       folder,
		ApplyUpdates: true,
	}
	lister := process.NewList(local, reqs)

	return &Folder{
		lister: lister,
	}, nil
}

func (f *Folder) Do() string {
	res, err := f.lister.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if len(res.Tasks) == 0 {
		return "no tasks here\n\n"
	}

	return format.FormatTaskTable(res.Tasks)
}
