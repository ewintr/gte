package task_test

import (
	"sort"
	"testing"
	"time"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/task"
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
	t.Run("no date", func(t *testing.T) {
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
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, task.NewDateFromString(tc.input))
			})
		}
	})

	t.Run("digits", func(t *testing.T) {
		task.Today = task.NewDate(2021, 1, 30)
		for _, tc := range []struct {
			name  string
			input string
			exp   task.Date
		}{
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
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, task.NewDateFromString(tc.input))

			})
		}
	})

	t.Run("day name", func(t *testing.T) {
		task.Today = task.NewDate(2021, 1, 30)
		for _, tc := range []struct {
			name  string
			input string
			exp   task.Date
		}{
			{
				name:  "dayname lowercase",
				input: "monday",
				exp:   task.NewDate(2021, 2, 1),
			},
			{
				name:  "dayname capitalized",
				input: "Monday",
				exp:   task.NewDate(2021, 2, 1),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, task.NewDateFromString(tc.input))
			})
		}
	})

	t.Run("relative days", func(t *testing.T) {
		task.Today = task.NewDate(2021, 1, 30)
		for _, tc := range []struct {
			name string
			exp  task.Date
		}{
			{
				name: "today",
				exp:  task.NewDate(2021, 1, 30),
			},
			{
				name: "tomorrow",
				exp:  task.NewDate(2021, 1, 31),
			},
			{
				name: "day after tomorrow",
				exp:  task.NewDate(2021, 2, 1),
			},
			{
				name: "this week",
				exp:  task.NewDate(2021, 2, 5),
			},
			{
				name: "next week",
				exp:  task.NewDate(2021, 2, 12),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, task.NewDateFromString(tc.name))
			})
		}
	})

	t.Run("sprint", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			today task.Date
			input string
			exp   task.Date
		}{
			{
				name:  "this sprint",
				today: task.NewDate(2021, 1, 30),
				input: "this sprint",
				exp:   task.NewDate(2021, 2, 11),
			},
			{
				name:  "jump week",
				today: task.NewDate(2021, 2, 5),
				input: "this sprint",
				exp:   task.NewDate(2021, 2, 11),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				task.Today = tc.today
				test.Equals(t, tc.exp, task.NewDateFromString(tc.input))
			})
		}
	})

	t.Run("empty", func(t *testing.T) {
		test.Equals(t, task.Date{}, task.NewDateFromString("test"))
	})
}

func TestDateDaysBetween(t *testing.T) {
	for _, tc := range []struct {
		name string
		d1   task.Date
		d2   task.Date
		exp  int
	}{
		{
			name: "same",
			d1:   task.NewDate(2021, 6, 23),
			d2:   task.NewDate(2021, 6, 23),
		},
		{
			name: "one",
			d1:   task.NewDate(2021, 6, 23),
			d2:   task.NewDate(2021, 6, 24),
			exp:  1,
		},
		{
			name: "many",
			d1:   task.NewDate(2021, 6, 23),
			d2:   task.NewDate(2024, 3, 7),
			exp:  988,
		},
		{
			name: "edge",
			d1:   task.NewDate(2020, 12, 30),
			d2:   task.NewDate(2021, 1, 3),
			exp:  4,
		},
		{
			name: "reverse",
			d1:   task.NewDate(2021, 6, 23),
			d2:   task.NewDate(2021, 5, 23),
			exp:  31,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.d1.DaysBetween(tc.d2))
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
			date: task.NewDate(2021, 5, 30),
			exp:  "2021-05-30 (sunday)",
		},
		{
			name: "normalize",
			date: task.NewDate(2021, 5, 32),
			exp:  "2021-06-01 (tuesday)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.date.String())
		})
	}
}

func TestDateHuman(t *testing.T) {
	monday := task.Today.Add(1)
	for {
		if monday.Weekday() == time.Monday {
			break
		}
		monday = monday.Add(1)
	}

	for _, tc := range []struct {
		name string
		date task.Date
		exp  string
	}{
		{
			name: "zero",
			date: task.NewDate(0, 0, 0),
			exp:  "-",
		},
		{
			name: "weekday",
			date: monday,
			exp:  "monday",
		},
		{
			name: "default",
			date: task.NewDate(2020, 1, 1),
			exp:  "2020-01-01 (wednesday)",
		},
		{
			name: "today",
			date: task.Today,
			exp:  "today",
		},
		{
			name: "tomorrow",
			date: task.Today.Add(1),
			exp:  "tomorrow",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.date.Human())
		})
	}
}

func TestDateIsZero(t *testing.T) {
	test.Equals(t, true, task.Date{}.IsZero())
	test.Equals(t, false, task.NewDate(2021, 6, 24).IsZero())
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
