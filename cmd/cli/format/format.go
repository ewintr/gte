package format

import (
	"fmt"
	"sort"

	"git.ewintr.nl/gte/internal/task"
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTaskTable(tasks []*task.LocalTask) string {
	if len(tasks) == 0 {
		return "no tasks to display\n"
	}

	sort.Sort(task.ByDefault(tasks))

	var output string
	for _, t := range tasks {
		var updateStr string
		if t.LocalStatus == task.STATUS_UPDATED {
			updateStr = " *"
		}
		output += fmt.Sprintf("%d%s\t%s\t%s (%s)\n", t.LocalId, updateStr, t.Due.String(), t.Action, t.Project)
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
