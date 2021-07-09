package command

import (
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

// New sends an action to the NEW folder so it can be updated to a real task later
type New struct {
	disp   *storage.Dispatcher
	action string
}

func NewNew(conf *configuration.Configuration, cmdArgs []string) (*New, error) {
	if len(cmdArgs) != 1 {
		return &New{}, ErrInvalidAmountOfArgs
	}

	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	return &New{
		disp:   disp,
		action: cmdArgs[0],
	}, nil
}

func (n *New) Do() string {
	if err := n.disp.Dispatch(&task.Task{Action: n.action}); err != nil {
		return FormatError(err)
	}

	return "message sent\n"
}
