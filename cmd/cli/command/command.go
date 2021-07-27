package command

import (
	"errors"
	"strconv"

	"git.ewintr.nl/gte/internal/configuration"
)

var (
	ErrInvalidAmountOfArgs = errors.New("invalid amount of args")
)

type Command interface {
	Do() string
	Cmd() string
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
	case "new":
		return NewNew(conf, cmdArgs)
	case "done":
		return NewDone(conf, cmdArgs)
	default:
		return NewEmpty()
	}
}

func parseTaskCommand(id int, tArgs []string, conf *configuration.Configuration) (Command, error) {
	if len(tArgs) == 0 {
		return NewShow(id, conf)
	}

	return NewEmpty()
}
