package format

import (
	"fmt"

	"code.ewintr.nl/gte/internal/task"
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
	tl := t
	tl.ApplyUpdate()

	output := fmt.Sprintf(`folder:  %s
status:  %s
action:  %s
project: %s
`, tl.Folder, tl.LocalStatus, tl.Action, tl.Project)
	if t.IsRecurrer() {
		output += fmt.Sprintf("recur:   %s\n", tl.Recur.String())
	} else {
		output += fmt.Sprintf("due:     %s\n", tl.Due.String())
	}

	return fmt.Sprintf("\n%s\n", output)
}
