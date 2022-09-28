package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/task"
)

var (
	ErrInvalidAmountOfArgs = errors.New("invalid amount of args")
	ErrInvalidArg          = errors.New("invalid argument")
	ErrCouldNotFindTask    = errors.New("could not find task")
	ErrUnknownFolder       = errors.New("unknown folder")
	ErrFieldAlreadyUsed    = errors.New("field was already used")
	ErrInvalidDate         = errors.New("could not understand date format")
	ErrInvalidProject      = errors.New("could not understand project")
)

type Command interface {
	Do() string
}

func Parse(args []string, conf *configuration.Configuration) (Command, error) {
	if len(args) == 0 {
		return NewEmpty()
	}
	cmd, cmdArgs := args[0], args[1:]

	id, err := strconv.Atoi(cmd)
	if err == nil {
		return parseTaskCommand(id, cmdArgs, conf)
	}

	switch cmd {
	case "fetch":
		return NewFetch(conf)
	case "send":
		return NewSend(conf)
	case "sync":
		return NewSync(conf)
	case "today":
		return NewToday(conf)
	case "tomorrow":
		return NewTomorrow(conf)
	case "week":
		return NewWeek(conf)
	case "project":
		return NewProject(conf, cmdArgs)
	case "projects":
		return NewProjects(conf)
	case "folder":
		return NewFolder(conf, cmdArgs)
	case "new":
		return NewNew(conf, cmdArgs)
	case "remote":
		return parseRemote(conf, cmdArgs)
	default:
		return NewEmpty()
	}
}

func parseTaskCommand(id int, tArgs []string, conf *configuration.Configuration) (Command, error) {
	if len(tArgs) == 0 {
		return NewShow(id, conf)
	}

	cmd, cmdArgs := tArgs[0], tArgs[1:]
	switch cmd {
	case "done":
		fallthrough
	case "del":
		return NewDone(id, conf)
	case "mod":
		return NewUpdate(id, conf, cmdArgs)
	default:
		return NewShow(id, conf)
	}
}

func parseRemote(conf *configuration.Configuration, cmdArgs []string) (Command, error) {
	switch {
	case len(cmdArgs) < 1:
		cmd, _ := NewEmpty()
		return cmd, ErrInvalidAmountOfArgs
	case cmdArgs[0] == "recur":
		return NewRecur(conf, cmdArgs[1:])
	case cmdArgs[0] == "inbox":
		return NewInbox(conf)
	default:
		cmd, _ := NewEmpty()
		return cmd, ErrInvalidArg
	}
}

func ParseTaskFieldArgs(args []string) (*task.LocalUpdate, error) {
	lu := &task.LocalUpdate{}

	action, fields := []string{}, []string{}
	for _, f := range args {
		if project, ok := parseProjectField(f); ok {
			if lu.Project != "" {
				return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_PROJECT)

			}
			if project == "" {
				return &task.LocalUpdate{}, ErrInvalidProject
			}
			lu.Project = project
			fields = append(fields, task.FIELD_PROJECT)
			continue
		}
		if due, ok := parseDueField(f); ok {
			if due.IsZero() {
				return &task.LocalUpdate{}, ErrInvalidDate
			}
			if !lu.Due.IsZero() {
				return &task.LocalUpdate{}, fmt.Errorf("%w: %s", ErrFieldAlreadyUsed, task.FIELD_DUE)
			}
			lu.Due = due
			fields = append(fields, task.FIELD_DUE)
			continue
		}
		if len(f) > 0 {
			action = append(action, f)
		}
	}

	if len(action) > 0 {
		lu.Action = strings.Join(action, " ")
		fields = append(fields, task.FIELD_ACTION)
	}

	lu.Fields = fields

	return lu, nil
}

func parseProjectField(s string) (string, bool) {
	if !strings.HasPrefix(s, "project:") && !strings.HasPrefix(s, "p:") {
		return "", false
	}
	split := strings.SplitN(s, ":", 2)

	return split[1], true
}

func parseDueField(s string) (task.Date, bool) {
	if !strings.HasPrefix(s, "due:") && !strings.HasPrefix(s, "d:") {
		return task.Date{}, false
	}
	split := strings.SplitN(s, ":", 2)

	due := task.NewDateFromString(split[1])

	return due, true
}
