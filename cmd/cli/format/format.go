package format

import (
	"fmt"

	"ewintr.nl/gte/internal/task"
)

type Column int

const (
	COL_ID Column = iota
	COL_STATUS
	COL_DUE
	COL_ACTION
	COL_PROJECT
)

var (
	COL_ALL = []Column{COL_ID, COL_STATUS, COL_DUE, COL_ACTION, COL_PROJECT}
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTask(t *task.LocalTask) string {
	output := fmt.Sprintf(`folder:  %s
action:  %s
project: %s
`, t.Folder, t.Action, t.Project)
	if t.IsRecurrer() {
		output += fmt.Sprintf("recur:   %s\n", t.Recur.String())
	} else {
		output += fmt.Sprintf("due:     %s\n", t.Due.String())
	}

	return fmt.Sprintf("\n%s\n", output)
}
