package task

import (
	"strings"
	"time"
)

type Period int
type Recurrer interface {
	RecursOn(date Date) bool
	FirstAfter(date Date) Date
	String() string
}

func NewRecurrer(recurStr string) Recurrer {
	terms := strings.Split(recurStr, ", ")
	if len(terms) < 3 {
		return nil
	}

	startDate, err := time.Parse("2006-01-02", terms[0])
	if err != nil {
		return nil
	}

	if terms[1] != "weekly" {
		return nil
	}

	if terms[2] != "wednesday" {
		return nil
	}

	year, month, date := startDate.Date()
	return Weekly{
		Start:   NewDate(year, int(month), date),
		Weekday: time.Wednesday,
	}
}

// yyyy-mm-dd, weekly, wednesday
type Weekly struct {
	Start   Date
	Weekday time.Weekday
}

func (w Weekly) RecursOn(date Date) bool {
	if !w.Start.After(date) {
		return false
	}

	return w.Weekday == date.Weekday()
}

func (w Weekly) FirstAfter(date Date) Date {
	//sd := w.Start.Weekday()

	return date
}

func (w Weekly) String() string {
	return "2021-01-31, weekly, wednesday"
}

/*
type BiWeekly struct {
	Start   Date
	Weekday Weekday
}

type RecurringTask struct {
	Action   string
	Start    Date
	Recurrer Recurrer
}
*/
