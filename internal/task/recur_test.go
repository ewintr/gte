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

func TestWeekly(t *testing.T) {
	weekly := task.Weekly{
		Start:   task.NewDate(2021, 1, 31), // a sunday
		Weekday: time.Wednesday,
	}
	weeklyStr := "2021-01-31 (sunday), weekly, wednesday"

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
				name: "wrong weekday",
				date: task.NewDate(2021, 2, 1), // a monday
			},
			{
				name: "right day",
				date: task.NewDate(2021, 2, 3), // a wednesday
				exp:  true,
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
