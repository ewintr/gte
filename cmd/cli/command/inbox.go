package command

import (
	"fmt"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/pkg/mstore"
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
