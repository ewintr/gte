package format

import (
	"fmt"

	"git.ewintr.nl/gte/internal/task"
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTaskTable(tasks []*task.Task) string {
	var output string
	for _, t := range tasks {
		output += fmt.Sprintf("%s\t%s\t%s\n", t.Id, t.Due.String(), t.Action)
	}

	return output
}
