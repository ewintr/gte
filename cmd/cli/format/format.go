package format

import (
	"errors"
	"fmt"
	"strings"

	"ewintr.nl/gte/internal/task"
)

var (
	ErrFieldAlreadyUsed = errors.New("field was already used")
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
	output := fmt.Sprintf(`folder: %s
action:  %s
project: %s
due:     %s
`, t.Folder, t.Action, t.Project, t.Due.String())
	if t.IsRecurrer() {
		output += fmt.Sprintf("recur:%s", t.Recur.String())
	}

	return fmt.Sprintf("%s\n", output)
}

func ParseTaskFieldArgs(args []string) (*task.LocalUpdate, error) {
	lu := &task.LocalUpdate{}

	action, fields := []string{}, []string{}
	for _, f := range args {
		split := strings.SplitN(f, ":", 2)
		if len(split) == 2 {
			switch split[0] {
			case "project":
				if lu.Project != "" {
					return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_PROJECT)
				}
				lu.Project = split[1]
				fields = append(fields, task.FIELD_PROJECT)
			case "due":
				if !lu.Due.IsZero() {
					return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_DUE)
				}
				lu.Due = task.NewDateFromString(split[1])
				fields = append(fields, task.FIELD_DUE)
			}
		} else {
			if len(f) > 0 {
				action = append(action, f)
			}
		}
	}

	if len(action) > 0 {
		lu.Action = strings.Join(action, " ")
		fields = append(fields, task.FIELD_ACTION)
	}

	lu.Fields = fields

	return lu, nil
}
