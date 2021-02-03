package task_test

import (
	"testing"
	"time"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
)

func TestDaily(t *testing.T) {
	daily := task.Daily{
		Start: task.NewDate(2021, 1, 31), // a sunday
	}
	dailyStr := "2021-01-31 (sunday), daily"

	t.Run("parse", func(t *testing.T) {
		test.Equals(t, daily, task.NewRecurrer(dailyStr))
	})

	t.Run("string", func(t *testing.T) {
		test.Equals(t, dailyStr, daily.String())
	})

	t.Run("recurs_on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date task.Date
			exp  bool
		}{
			{
				name: "before",
				date: task.NewDate(2021, 1, 30),
			},
			{
				name: "on",
				date: daily.Start,
				exp:  true,
			},
			{
				name: "after",
				date: task.NewDate(2021, 2, 1),
				exp:  true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, daily.RecursOn(tc.date))
			})
		}
	})
}

func TestParseWeekly(t *testing.T) {
	start := task.NewDate(2021, 2, 7)
	for _, tc := range []struct {
		name      string
		input     []string
		expOk     bool
		expWeekly task.Weekly
	}{
		{
			name: "empty",
		},
		{
			name:  "wrong type",
			input: []string{"daily"},
		},
		{
			name:  "wrong count",
			input: []string{"weeekly"},
		},
		{
			name:  "unknown day",
			input: []string{"weekly", "festivus"},
		},
		{
			name:  "one day",
			input: []string{"weekly", "monday"},
			expOk: true,
			expWeekly: task.Weekly{
				Start: start,
				Weekdays: task.Weekdays{
					time.Monday,
				},
			},
		},
		{
			name:  "multiple days",
			input: []string{"weekly", "monday & thursday & saturday"},
			expOk: true,
			expWeekly: task.Weekly{
				Start: start,
				Weekdays: task.Weekdays{
					time.Monday,
					time.Thursday,
					time.Saturday,
				},
			},
		},
		{
			name:  "wrong order",
			input: []string{"weekly", "sunday & thursday & wednesday"},
			expOk: true,
			expWeekly: task.Weekly{
				Start: start,
				Weekdays: task.Weekdays{
					time.Wednesday,
					time.Thursday,
					time.Sunday,
				},
			},
		},
		{
			name:  "doubles",
			input: []string{"weekly", "sunday & sunday & monday"},
			expOk: true,
			expWeekly: task.Weekly{
				Start: start,
				Weekdays: task.Weekdays{
					time.Monday,
					time.Sunday,
				},
			},
		},
		{
			name:  "one unknown",
			input: []string{"weekly", "sunday & someday"},
			expOk: true,
			expWeekly: task.Weekly{
				Start: start,
				Weekdays: task.Weekdays{
					time.Sunday,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			weekly, ok := task.ParseWeekly(start, tc.input)
			test.Equals(t, tc.expOk, ok)
			if tc.expOk {
				test.Equals(t, tc.expWeekly, weekly)
			}
		})
	}
}

func TestWeekly(t *testing.T) {
	weekly := task.Weekly{
		Start: task.NewDate(2021, 1, 31), // a sunday
		Weekdays: task.Weekdays{
			time.Monday,
			time.Wednesday,
			time.Thursday,
		},
	}
	weeklyStr := "2021-01-31 (sunday), weekly, monday & wednesday & thursday"

	t.Run("parse", func(t *testing.T) {
		test.Equals(t, weekly, task.NewRecurrer(weeklyStr))
	})

	t.Run("string", func(t *testing.T) {
		test.Equals(t, weeklyStr, weekly.String())
	})

	t.Run("recurs_on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date task.Date
			exp  bool
		}{
			{
				name: "before start",
				date: task.NewDate(2021, 1, 27), // a wednesday
			},
			{
				name: "right weekday",
				date: task.NewDate(2021, 2, 1), // a monday
				exp:  true,
			},
			{
				name: "another right day",
				date: task.NewDate(2021, 2, 3), // a wednesday
				exp:  true,
			},
			{
				name: "wrong weekday",
				date: task.NewDate(2021, 2, 5), // a friday
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, weekly.RecursOn(tc.date))
			})
		}
	})
}

func TestBiweekly(t *testing.T) {
	biweekly := task.Biweekly{
		Start:   task.NewDate(2021, 1, 31), // a sunday
		Weekday: time.Wednesday,
	}
	biweeklyStr := "2021-01-31 (sunday), biweekly, wednesday"

	t.Run("parse", func(t *testing.T) {
		test.Equals(t, biweekly, task.NewRecurrer(biweeklyStr))
	})

	t.Run("string", func(t *testing.T) {
		test.Equals(t, biweeklyStr, biweekly.String())
	})

	t.Run("recurs_on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date task.Date
			exp  bool
		}{
			{
				name: "before start",
				date: task.NewDate(2021, 1, 27), // a wednesday
			},
			{
				name: "wrong weekday",
				date: task.NewDate(2021, 2, 1), // a monday
			},
			{
				name: "odd week count",
				date: task.NewDate(2021, 2, 10), // a wednesday
			},
			{
				name: "right",
				date: task.NewDate(2021, 2, 17), // a wednesday
				exp:  true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, biweekly.RecursOn(tc.date))
			})
		}
	})
}

func TestEveryNWeeks(t *testing.T) {
	everyNWeeks := task.EveryNWeeks{
		Start: task.NewDate(2021, 2, 3),
		N:     3,
	}
	everyNWeeksStr := "2021-02-03 (wednesday), every 3 weeks"

	t.Run("parse", func(t *testing.T) {
		test.Equals(t, everyNWeeks, task.NewRecurrer(everyNWeeksStr))
	})

	t.Run("string", func(t *testing.T) {
		test.Equals(t, everyNWeeksStr, everyNWeeks.String())
	})

	t.Run("recurs on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date task.Date
			exp  bool
		}{
			{
				name: "before start",
				date: task.NewDate(2021, 1, 27),
			},
			{
				name: "on start",
				date: task.NewDate(2021, 2, 3),
				exp:  true,
			},
			{
				name: "wrong day",
				date: task.NewDate(2021, 2, 4),
			},
			{
				name: "one week after",
				date: task.NewDate(2021, 2, 10),
			},
			{
				name: "first interval",
				date: task.NewDate(2021, 2, 24),
				exp:  true,
			},
			{
				name: "second interval",
				date: task.NewDate(2021, 3, 17),
				exp:  true,
			},
			{
				name: "second interval plus one week",
				date: task.NewDate(2021, 3, 24),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				test.Equals(t, tc.exp, everyNWeeks.RecursOn(tc.date))
			})
		}
	})
}
