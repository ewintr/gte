package task_test

import (
	"testing"
	"time"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
)

func TestNewRecurrer(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		exp   task.Recurrer
	}{
		{
			name: "empty",
		},
		{
			name:  "weekly",
			input: "2021-01-31, weekly, wednesday",
			exp: task.Weekly{
				Start:   task.NewDate(2021, 1, 31),
				Weekday: time.Wednesday,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, task.NewRecurrer(tc.input))
		})
	}
}
