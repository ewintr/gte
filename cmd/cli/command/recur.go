package command

import (
	"fmt"
	"strconv"

	"go-mod.ewintr.nl/gte/cmd/cli/format"
	"go-mod.ewintr.nl/gte/internal/configuration"
	"go-mod.ewintr.nl/gte/internal/process"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/pkg/msend"
	"go-mod.ewintr.nl/gte/pkg/mstore"
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

	return fmt.Sprintf("generated %d tasks\n\n", res.Count)
}
