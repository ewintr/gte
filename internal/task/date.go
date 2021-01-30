package task

import (
	"time"
)

type Date struct {
	t time.Time
}

func NewDate(year, month, day int) *Date {
	var m time.Month
	switch month {
	case 1:
		m = time.January
	case 2:
		m = time.February
	case 3:
		m = time.March
	case 4:
		m = time.April
	case 5:
		m = time.May
	case 6:
		m = time.June
	case 7:
		m = time.July
	case 8:
		m = time.August
	case 9:
		m = time.September
	case 10:
		m = time.October
	case 11:
		m = time.November
	case 12:
		m = time.December
	}

	t := time.Date(year, m, day, 10, 0, 0, 0, time.UTC)

	if year == 0 && month == 0 && day == 0 {
		t = time.Time{}
	}

	return &Date{
		t: t,
	}
}

func (d *Date) String() string {
	if d.t.IsZero() {
		return "no date"
	}

	return d.t.Format("2006-01-02")
}
