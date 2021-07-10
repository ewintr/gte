package command

import (
	"fmt"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/pkg/mstore"
)

type Sync struct {
	syncer *process.Sync
}

func NewSync(conf *configuration.Configuration) (*Sync, error) {
	msgStore := mstore.NewIMAP(conf.IMAP())
	remote := storage.NewRemoteRepository(msgStore)
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Sync{}, err
	}
	syncer := process.NewSync(remote, local)

	return &Sync{
		syncer: syncer,
	}, nil
}

func (s *Sync) Do() string {
	result, err := s.syncer.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("synced %d tasks\n", result.Count)
}
