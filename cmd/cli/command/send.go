package command

import (
	"fmt"

	"git.ewintr.nl/gte/cmd/cli/format"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/pkg/msend"
)

type Send struct {
	sender *process.Send
}

func NewSend(conf *configuration.Configuration) (*Send, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Send{}, err
	}
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	return &Send{
		sender: process.NewSend(local, disp),
	}, nil
}

func (s *Send) Do() string {
	count, err := s.sender.Process()
	if err != nil {
		return format.FormatError(err)
	}

	return fmt.Sprintf("sent %d tasks\n\n", count)
}
