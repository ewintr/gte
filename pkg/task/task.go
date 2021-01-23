package task

import "time"

type Date time.Time

func (d *Date) Weekday() Weekday {
	return d.Weekday()
}

type Weekday time.Weekday

type Period int

type Task struct {
	Action string
	Due    Date
}

type Recurrer interface {
	FirstAfter(date Date) Date
}

type Weekly struct {
	Start   Date
	Weekday Weekday
}

func (w *Weekly) FirstAfter(date Date) Date {
	//sd := w.Start.Weekday()

	return date
}

type BiWeekly struct {
	Start   Date
	Weekday Weekday
}

type RecurringTask struct {
	Action   string
	Start    Date
	Recurrer Recurrer
}
