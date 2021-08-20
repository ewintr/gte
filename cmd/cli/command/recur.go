package command

import (
	"fmt"
	"strconv"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

type Recur struct {
	recurrer *process.Recur
}

func NewRecur(conf *configuration.Configuration, cmdArgs []string) (*Recur, error) {
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	if len(cmdArgs) < 1 {
		return &Recur{}, ErrInvalidAmountOfArgs
	}
	daysAhead, err := strconv.Atoi(cmdArgs[0])
	if err != nil {
		return &Recur{}, ErrInvalidArg
	}

	return &Recur{
		recurrer: process.NewRecur(remote, disp, daysAhead),
	}, nil
}

func (r *Recur) Do() string {
	res, err := r.recurrer.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("generated %d tasks\n", res.Count)
}