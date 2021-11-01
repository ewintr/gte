package command

import (
	"fmt"
	"time"

	"ewintr.nl/gte/cmd/cli/format"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/pkg/msend"
	"ewintr.nl/gte/pkg/mstore"
)

type Sync struct {
	fetcher       *process.Fetch
	sender        *process.Send
	fetchInterval time.Duration
	fetchLatest   time.Time
}

func NewSync(conf *configuration.Configuration) (*Sync, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Sync{}, err
	}
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	fetchLatest, err := local.LatestSync()
	if err != nil {
		return &Sync{}, err
	}
	fetchInterval := 15 * time.Minute // not yet configurable

	return &Sync{
		fetcher:       process.NewFetch(remote, local),
		sender:        process.NewSend(local, disp),
		fetchInterval: fetchInterval,
		fetchLatest:   fetchLatest,
	}, nil
}

func (s *Sync) Do() string {
	countSend, err := s.sender.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if time.Now().Before(s.fetchLatest.Add(s.fetchInterval)) {
		return fmt.Sprintf("sent %d tasks, not time to fetch yet\n", countSend)
	}

	fResult, err := s.fetcher.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("sent %d, fetched %d tasks\n", countSend, fResult.Count)
}
