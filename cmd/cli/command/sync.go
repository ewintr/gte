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
	dispInterval  time.Duration
	dispLatest    time.Time
}

func NewSync(conf *configuration.Configuration) (*Sync, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Sync{}, err
	}
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	fetchLatest, dispLatest, err := local.LatestSyncs()
	if err != nil {
		return &Sync{}, err
	}
	fetchInterval := 15 * time.Minute // not yet configurable
	dispInterval := 2 * time.Minute

	return &Sync{
		fetcher:       process.NewFetch(remote, local),
		sender:        process.NewSend(local, disp),
		fetchInterval: fetchInterval,
		fetchLatest:   fetchLatest,
		dispInterval:  dispInterval,
		dispLatest:    dispLatest,
	}, nil
}

func (s *Sync) Do() string {
	countSend, err := s.sender.Process()
	if err != nil {
		return format.FormatError(err)
	}
	if countSend > 0 {
		return fmt.Sprintf("sent %d tasks, not fetching yet\n", countSend)
	}

	if time.Now().Before(s.dispLatest.Add(s.dispInterval)) {
		return "sent 0 tasks, send interval has not passed yet\n"
	}

	if time.Now().Before(s.fetchLatest.Add(s.fetchInterval)) {
		return "sent 0 tasks, fetch interval has not passed yet\n"
	}

	fResult, err := s.fetcher.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("fetched %d tasks\n", fResult.Count)
}
