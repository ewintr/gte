package command

import (
	"strings"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

// Add sends an action to the NEW folder so it can be updated to a real task later
type Add struct {
	disp   *storage.Dispatcher
	action string
}

func NewAdd(conf *configuration.Configuration, cmdArgs []string) (*Add, error) {
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	return &Add{
		disp:   disp,
		action: strings.Join(cmdArgs, " "),
	}, nil
}

func (n *Add) Do() string {
	if err := n.disp.Dispatch(&task.Task{Action: n.action}); err != nil {
		return format.FormatError(err)
	}

	return "message sent\n"
}
