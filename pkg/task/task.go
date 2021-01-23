package task

import "time"

type Date struct {
	date time.Time
}

type Weekday int

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
