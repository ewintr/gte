package command_test

import (
	"errors"
	"strings"
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/cmd/cli/command"
	"git.ewintr.nl/gte/internal/task"
)

func TestParseTaskFieldArgs(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		expUpdate task.LocalUpdate
		expErr    error
	}{
		{
			name:      "empty",
			expUpdate: task.LocalUpdate{},
		},
		{
			name:  "join action",
			input: "some things to do",
			expUpdate: task.LocalUpdate{
				Action: "some things to do",
			},
		},
		{
			name:  "all",
			input: "project:project do stuff due:2021-08-06",
			expUpdate: task.LocalUpdate{
				Action:  "do stuff",
				Project: "project",
				Due:     task.NewDate(2021, 8, 6),
			},
		},
		{
			name:  "no action",
			input: "due:2021-08-06",
			expUpdate: task.LocalUpdate{
				Due: task.NewDate(2021, 8, 6),
			},
		},
		{
			name:      "two projects",
			input:     "project:project1 project:project2",
			expUpdate: task.LocalUpdate{},
			expErr:    command.ErrFieldAlreadyUsed,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := strings.Split(tc.input, " ")
			act, err := command.ParseTaskFieldArgs(args)
			test.Equals(t, tc.expUpdate, act)
			test.Assert(t, errors.Is(err, tc.expErr), "wrong err")
		})
	}
}
