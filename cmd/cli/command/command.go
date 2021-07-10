package command

import (
	"errors"

	"git.ewintr.nl/gte/internal/configuration"
)

var (
	ErrInvalidAmountOfArgs = errors.New("invalid amount of args")
)

type Command interface {
	Do() string
}

func Parse(args []string, conf *configuration.Configuration) (Command, error) {
	if len(args) == 0 {
		return NewEmpty()
	}

	cmd, cmdArgs := args[0], args[1:]
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
