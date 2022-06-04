package format

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"ewintr.nl/gte/internal/task"
)

var (
	ErrFieldAlreadyUsed = errors.New("field was already used")
)

type column int

const (
	COL_ID column = iota
	COL_STATUS
	COL_DATE
	COL_TASK
	COL_PROJECT
)

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}

func FormatTasks(tasks []*task.LocalTask) string {
	if len(tasks) == 0 {
		return "no tasks to display\n"
	}

	sort.Sort(task.ByDefault(tasks))

	normal, recurring := []*task.LocalTask{}, []*task.LocalTask{}
	for _, lt := range tasks {
		if lt.Folder == task.FOLDER_RECURRING {
			recurring = append(recurring, lt)
			continue
		}
		normal = append(normal, lt)
	}

	output := ""
	if len(recurring) > 0 {
		output = fmt.Sprintf("%s\nRecurring:\n%s\nTasks:\n", output, FormatTaskTable(recurring))
	}

	return fmt.Sprintf("%s\n%s\n", output, FormatTaskTable(normal))
}

func FormatTaskTable(tasks []*task.LocalTask, cols ...column) string {
	if len(cols) == 0 {
		cols = []column{COL_ID, COL_STATUS, COL_DATE, COL_TASK, COL_PROJECT}
	}

	var data [][]string
	for _, t := range tasks {
		var line []string
		for _, col := range cols {
			switch col {
			case COL_ID:
				line = append(line, fmt.Sprintf("%d", t.LocalId))
			case COL_STATUS:
				var updated string
				if t.LocalStatus == task.STATUS_UPDATED {
					updated = "*"
				}
				line = append(line, updated)
			case COL_DATE:
				line = append(line, t.Due.String())
			case COL_TASK:
				line = append(line, t.Action)
			case COL_PROJECT:
				line = append(line, t.Project)
			}
		}
		data = append(data, line)
	}

	return FormatTable(data)
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

func FormatTable(data [][]string) string {
	if len(data) == 0 {
		return ""
	}
	max := make([]int, len(data))
	for _, line := range data {
		for i, col := range line {
			if len(col) > max[i] {
				max[i] = len(col)
			}
		}
	}

	var output string
	for _, line := range data {
		for i, col := range line {
			output += col
			for s := 0; s < max[i]-len(col); s++ {
				output += " "
			}
			output += " "
		}
		output += "\n"

	}

	return output
}
