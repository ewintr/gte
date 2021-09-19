package command

import (
	"fmt"

	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/pkg/mstore"
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
