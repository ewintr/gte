package command

import (
	"fmt"

	"go-mod.ewintr.nl/gte/cmd/cli/format"
	"go-mod.ewintr.nl/gte/internal/configuration"
	"go-mod.ewintr.nl/gte/internal/process"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/pkg/mstore"
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
