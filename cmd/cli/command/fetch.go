package command

import (
	"fmt"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/pkg/mstore"
)

type Fetch struct {
	fetcher *process.Fetch
}

func NewFetch(conf *configuration.Configuration) (*Fetch, error) {
	msgStore := mstore.NewIMAP(conf.IMAP())
	remote := storage.NewRemoteRepository(msgStore)
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Fetch{}, err
	}
	fetcher := process.NewFetch(remote, local)

	return &Fetch{
		fetcher: fetcher,
	}, nil
}

func (s *Fetch) Do() string {
	result, err := s.fetcher.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("fetched %d tasks\n\n", result.Count)
}
