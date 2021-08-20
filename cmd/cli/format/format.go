package format

import (
	"fmt"
	"sort"
	"strings"

	"git.ewintr.nl/gte/internal/task"
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTaskTable(tasks []*task.LocalTask) string {
	if len(tasks) == 0 {
		return "no tasks to display\n"
	}

	sort.Sort(task.ByDue(tasks))

	var output string
	for _, t := range tasks {
		output += fmt.Sprintf("%d\t%s\t%s (%s, %s)\n", t.LocalId, t.Due.String(), t.Action, t.Project, strings.ToLower(t.Folder))
	}

	return output
}

func FormatTask(t *task.LocalTask) string {
	output := fmt.Sprintf(`folder: %s
action:  %s
project: %s
due:     %s
`, t.Folder, t.Action, t.Project, t.Due.String())
	if t.IsRecurrer() {
		output += fmt.Sprintf("recur:%s", t.Recur.String())
	}

	return output
}
