package command_test

import (
	"errors"
	"strings"
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/cmd/cli/command"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
)

func TestParseTaskFieldArgs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expField process.UpdateFields
		expErr   error
	}{
		{
			name: "empty",
			expField: process.UpdateFields{
				task.FIELD_ACTION: "",
			},
		},
		{
			name:  "join action",
			input: "some things to do",
			expField: process.UpdateFields{
				task.FIELD_ACTION: "some things to do",
			},
		},
		{
			name:  "all",
			input: "project:project do stuff due:2021-08-06",
			expField: process.UpdateFields{
				task.FIELD_ACTION:  "do stuff",
				task.FIELD_PROJECT: "project",
				task.FIELD_DUE:     "2021-08-06",
			},
		},
		{
			name:  "no action",
			input: "due:2021-08-06",
			expField: process.UpdateFields{
				task.FIELD_DUE: "2021-08-06",
			},
		},
		{
			name:     "two projects",
			input:    "project:project1 project:project2",
			expField: process.UpdateFields{},
			expErr:   command.ErrFieldAlreadyUsed,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := strings.Split(tc.input, " ")
			act, err := command.ParseTaskFieldArgs(args)
			test.Equals(t, tc.expField, act)
			test.Assert(t, errors.Is(err, tc.expErr), "wrong err")
		})
	}
}
