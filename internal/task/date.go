package task

import "time"

type Date time.Time

func (d *Date) Weekday() Weekday {
	return d.Weekday()
}
