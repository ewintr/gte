package task_test

import (
	"sort"
	"testing"
	"time"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
)

func TestWeekdaysSort(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input task.Weekdays
		exp   task.Weekdays
	}{
		{
			name: "empty",
		},
		{
			name:  "one",
			input: task.Weekdays{time.Tuesday},
			exp:   task.Weekdays{time.Tuesday},
		},
		{
			name:  "multiple",
			input: task.Weekdays{time.Wednesday, time.Tuesday, time.Monday},
			exp:   task.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
		},
		{
			name:  "sunday is last",
			input: task.Weekdays{time.Saturday, time.Sunday, time.Monday},
			exp:   task.Weekdays{time.Monday, time.Saturday, time.Sunday},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(tc.input)
			test.Equals(t, tc.exp, tc.input)
		})
	}
}

func TestWeekdaysUnique(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input task.Weekdays
		exp   task.Weekdays
	}{
		{
			name:  "empty",
			input: task.Weekdays{},
			exp:   task.Weekdays{},
		},
		{
			name:  "single",
			input: task.Weekdays{time.Monday},
			exp:   task.Weekdays{time.Monday},
		},
		{
			name:  "no doubles",
			input: task.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
			exp:   task.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
		},
		{
			name:  "doubles",
			input: task.Weekdays{time.Monday, time.Monday, time.Wednesday, time.Monday},
			exp:   task.Weekdays{time.Monday, time.Wednesday},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.input.Unique())
		})
	}
}

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
			name:  "dutch dayname lowercase",
			input: "maandag",
			exp:   task.NewDate(2021, 2, 1),
		},
		{
			name:  "today english",
			input: "today",
			exp:   task.NewDate(2021, 1, 30),
		},
		{
			name:  "today dutch",
			input: "vandaag",
			exp:   task.NewDate(2021, 1, 30),
		},
		{
			name:  "tomorrow english",
			input: "tomorrow",
			exp:   task.NewDate(2021, 1, 31),
		},
		{
			name:  "tomorrow dutch",
			input: "morgen",
			exp:   task.NewDate(2021, 1, 31),
		},
		{
			name:  "day after tomorrow english",
			input: "day after tomorrow",
			exp:   task.NewDate(2021, 2, 1),
		},
		{
			name:  "day after tomorrow dutch",
			input: "overmorgen",
			exp:   task.NewDate(2021, 2, 1),
		},
		{
			name:  "this week english",
			input: "this week",
			exp:   task.NewDate(2021, 2, 5),
		},
		{
			name:  "this week dutch",
			input: "deze week",
			exp:   task.NewDate(2021, 2, 5),
		},
		{
			name:  "next week english",
			input: "next week",
			exp:   task.NewDate(2021, 2, 12),
		},
		{
			name:  "next week dutch",
			input: "volgende week",
			exp:   task.NewDate(2021, 2, 12),
		},
		{
			name:  "this sprint english",
			input: "this sprint",
			exp:   task.NewDate(2021, 2, 11),
		},
		{
			name:  "this sprint dutch",
			input: "deze sprint",
			exp:   task.NewDate(2021, 2, 11),
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
