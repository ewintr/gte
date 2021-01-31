package task_test

import (
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
)

func TestNewDateFromString(t *testing.T) {
	task.Today = task.NewDate(2021, 1, 30)
	for _, tc := range []struct {
		name  string
		input string
		exp   task.Date
	}{
		{
			name: "empty",
			exp:  task.Date{},
		},
		{
			name:  "no date",
			input: "no date",
			exp:   task.Date{},
		},
		{
			name:  "normal",
			input: "2021-01-30 (saturday)",
			exp:   task.NewDate(2021, 1, 30),
		},
		{
			name:  "short",
			input: "2021-01-30",
			exp:   task.NewDate(2021, 1, 30),
		},
		{
			name:  "english dayname lowercase",
			input: "monday",
			exp:   task.NewDate(2021, 2, 1),
		},
		{
			name:  "english dayname capitalized",
			input: "Monday",
			exp:   task.NewDate(2021, 2, 1),
		},
		{
			name:  "ducth dayname lowercase",
			input: "maandag",
			exp:   task.NewDate(2021, 2, 1),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, task.NewDateFromString(tc.input))
		})
	}
}

func TestDateString(t *testing.T) {
	for _, tc := range []struct {
		name string
		date task.Date
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
			exp:  "2021-01-30 (saturday)",
		},
		{
			name: "normalize",
			date: task.NewDate(2021, 1, 32),
			exp:  "2021-02-01 (monday)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.date.String())
		})
	}
}

func TestDateAfter(t *testing.T) {
	day := task.NewDate(2021, 1, 31)
	for _, tc := range []struct {
		name string
		tDay task.Date
		exp  bool
	}{
		{
			name: "after",
			tDay: task.NewDate(2021, 1, 30),
			exp:  true,
		},
		{
			name: "on",
			tDay: day,
		},
		{
			name: "before",
			tDay: task.NewDate(2021, 2, 1),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, day.After(tc.tDay))
		})
	}
}
