package command

import (
	"sort"
	"strings"

	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
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

	sort.Sort(task.ByDefault(res.Tasks))

	return format.FormatTaskTable(res.Tasks, format.COL_ALL)
}
