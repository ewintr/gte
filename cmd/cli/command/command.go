package command

import (
	"fmt"

	"git.ewintr.nl/gte/internal/configuration"
)

type Command interface {
	Do() string
}

func Parse(args []string, conf *configuration.Configuration) (Command, error) {
	if len(args) == 0 {
		return NewEmpty()
	}

	cmd, _ := args[0], args[1:]
	switch cmd {
	case "sync":
		return NewSync(conf)
	case "today":
		return NewToday(conf)
	default:
		return NewEmpty()
	}
}

func FormatError(err error) string {
	return fmt.Sprintf("could not perform command.\n\nerror: %s\n", err.Error())
}
