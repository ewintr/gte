package command

import (
	"fmt"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/pkg/mstore"
)

type Inbox struct {
	inboxer *process.Inbox
}

func NewInbox(conf *configuration.Configuration) (*Inbox, error) {
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))

	return &Inbox{
		inboxer: process.NewInbox(remote),
	}, nil
}

func (i *Inbox) Do() string {
	res, err := i.inboxer.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("processed %d tasks\n\n", res.Count)
}
