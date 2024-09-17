package command_test

import (
	"errors"
	"strings"
	"testing"

	"go-mod.ewintr.nl/go-kit/test"
	"go-mod.ewintr.nl/gte/cmd/cli/command"
	"go-mod.ewintr.nl/gte/internal/task"
)

func TestParseTaskFieldArgs(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		expUpdate *task.LocalUpdate
		expErr    error
	}{
		{
			name: "empty",
			expUpdate: &task.LocalUpdate{
				Fields: []string{},
			},
		},
		{
			name:  "join action",
			input: "some things to do",
			expUpdate: &task.LocalUpdate{
				Fields: []string{task.FIELD_ACTION},
				Action: "some things to do",
			},
		},
		{
			name:  "all",
			input: "project:project do stuff due:2021-08-06",
			expUpdate: &task.LocalUpdate{
				Fields:  []string{task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_ACTION},
				Action:  "do stuff",
				Project: "project",
				Due:     task.NewDate(2021, 8, 6),
			},
		},
		{
			name:  "no action",
			input: "due:2021-08-06",
			expUpdate: &task.LocalUpdate{
				Fields: []string{task.FIELD_DUE},
				Due:    task.NewDate(2021, 8, 6),
			},
		},
		{
			name:      "two projects",
			input:     "project:project1 project:project2",
			expUpdate: &task.LocalUpdate{},
			expErr:    command.ErrFieldAlreadyUsed,
		},
		{
			name:  "abbreviated",
			input: "p:project1 d:2022-09-28",
			expUpdate: &task.LocalUpdate{
				Fields:  []string{task.FIELD_PROJECT, task.FIELD_DUE},
				Project: "project1",
				Due:     task.NewDate(2022, 9, 28),
			},
		},
		{
			name:      "empty project",
			input:     "action project:",
			expUpdate: &task.LocalUpdate{},
			expErr:    command.ErrInvalidProject,
		},
		{
			name:      "empty date",
			input:     "action due:",
			expUpdate: &task.LocalUpdate{},
			expErr:    command.ErrInvalidDate,
		},
		{
			name:  "url",
			input: "https://ewintr.nl/something?arg=1",
			expUpdate: &task.LocalUpdate{
				Fields: []string{task.FIELD_ACTION},
				Action: "https://ewintr.nl/something?arg=1",
			},
		},
		{
			name:      "misformatted date",
			input:     "d:20-wrong",
			expUpdate: &task.LocalUpdate{},
			expErr:    command.ErrInvalidDate,
		},
		{
			name:  "recur",
			input: "recur:today,daily",
			expUpdate: &task.LocalUpdate{
				Fields: []string{task.FIELD_RECUR},
				Recur:  task.NewRecurrer("today, daily"),
			},
		},
		{
			name:  "recurs short",
			input: "r:today,daily",
			expUpdate: &task.LocalUpdate{
				Fields: []string{task.FIELD_RECUR},
				Recur:  task.NewRecurrer("today, daily"),
			},
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
