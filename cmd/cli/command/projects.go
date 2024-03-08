package command

import (
	"fmt"

	"code.ewintr.nl/gte/cmd/cli/format"
	"code.ewintr.nl/gte/internal/configuration"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
)

type Projects struct {
	local     storage.LocalRepository
	projecter *process.Projects
}

func NewProjects(conf *configuration.Configuration) (*Projects, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Projects{}, err
	}

	projecter := process.NewProjects(local)

	return &Projects{
		local:     local,
		projecter: projecter,
	}, nil
}

func (p *Projects) Do() string {
	projects, err := p.projecter.Process()
	if err != nil {
		return format.FormatError(err)
	}

	if len(projects) == 0 {
		return "no projects here\n\n"
	}

	var data [][]string
	for _, project := range projects {
		if project != "" {
			data = append(data, []string{project})
		}
	}

	return fmt.Sprintf("%s", format.FormatTable(data))
}
