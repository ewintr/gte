package command

import (
	"errors"
	"strconv"

	"git.ewintr.nl/gte/internal/configuration"
)

var (
	ErrInvalidAmountOfArgs = errors.New("invalid amount of args")
	ErrCouldNotFindTask    = errors.New("could not find task")
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
	case "sync":
		return NewSync(conf)
	case "today":
		return NewToday(conf)
	case "tomorrow":
		return NewTomorrow(conf)
	case "project":
		return NewProject(conf, cmdArgs)
	case "projects":
		return NewProjects(conf)
	case "folder":
		return NewFolder(conf, cmdArgs)
	case "add":
		return NewAdd(conf, cmdArgs)
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
