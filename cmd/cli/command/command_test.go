package command_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/cmd/cli/command"
	"git.ewintr.nl/gte/internal/configuration"
)

func TestCommand(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		exp  string
	}{
		{
			name: "empty",
			exp:  "empty",
		},
		{
			name: "sync",
			args: []string{"sync"},
			exp:  "sync",
		},
		{
			name: "today",
			args: []string{"today"},
			exp:  "today",
		},
		{
			name: "tomorrow",
			args: []string{"tomorrow"},
			exp:  "tomorrow",
		},
		{
			name: "done",
			args: []string{"done"},
			exp:  "done",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _ := command.Parse(tc.args, &configuration.Configuration{})
			test.Equals(t, tc.exp, cmd.Cmd())
		})
	}
}
