package format

import (
	"fmt"

	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTaskTable(local storage.LocalRepository, tasks []*task.Task) string {
	if len(tasks) == 0 {
		return "no tasks to display\n"
	}

	localIds, err := local.LocalIds()
	if err != nil {
		return FormatError(err)
	}

	var output string
	for _, t := range tasks {
		output += fmt.Sprintf("%d\t%s\t%s\n", localIds[t.Id], t.Due.String(), t.Action)
	}

	return output
}
