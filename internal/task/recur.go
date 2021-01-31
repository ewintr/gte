package task

import "time"

type Weekday time.Weekday

type Period int
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
