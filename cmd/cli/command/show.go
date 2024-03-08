package command

import (
	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/storage"
)

type Show struct {
	local storage.LocalRepository
	id    int
}

func NewShow(id int, conf *configuration.Configuration) (*Show, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Show{}, err
	}

	return &Show{
		local: local,
		id:    id,
	}, nil
}

func (s *Show) Do() string {
	t, err := s.local.FindByLocalId(s.id)
	if err != nil {
		return format.FormatError(err)
	}

	return format.FormatTask(t)
}
