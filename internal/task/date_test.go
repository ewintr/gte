package task_test

import (
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
)

func TestDateString(t *testing.T) {
	for _, tc := range []struct {
		name string
		date *task.Date
		exp  string
	}{
		{
			name: "zero",
			date: task.NewDate(0, 0, 0),
			exp:  "no date",
		},
		{
			name: "normal",
			date: task.NewDate(2021, 1, 30),
			exp:  "2021-01-30",
		},
		{
			name: "normalize",
			date: task.NewDate(2021, 1, 32),
			exp:  "2021-02-01",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.date.String())
		})
	}
}
