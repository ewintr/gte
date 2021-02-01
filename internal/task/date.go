package task

import (
	"fmt"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02 (Monday)"
)

var Today Date

func init() {
	year, month, day := time.Now().Date()
	Today = NewDate(year, int(month), day)
}

type Date struct {
	t time.Time
}

func NewDate(year, month, day int) Date {

	if year == 0 && month == 0 && day == 0 {
		return Date{}
	}

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

	t := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)

	return Date{
		t: t,
	}
}

func NewDateFromString(date string) Date {
	date = strings.ToLower(strings.TrimSpace(date))

	if date == "no date" || date == "" {
		return Date{}
	}

	if date == "today" || date == "vandaag" {
		return Today
	}

	if date == "tomorrow" || date == "morgen" {
		return Today.AddDays(1)
	}

	t, err := time.Parse("2006-01-02", fmt.Sprintf("%.10s", date))
	if err == nil {
		return Date{t: t}
	}

	weekday := Today.Weekday()
	newWeekday, ok := ParseWeekday(date)
	if !ok {
		return Date{}
	}
	daysToAdd := int(newWeekday) - int(weekday)
	if daysToAdd <= 0 {
		daysToAdd += 7
	}

	return Today.Add(daysToAdd)
}

func (d *Date) String() string {
	if d.t.IsZero() {
		return "no date"
	}

	return strings.ToLower(d.t.Format(DateFormat))
}

func (d *Date) IsZero() bool {
	return d.t.IsZero()
}

func (d *Date) Time() time.Time {
	return d.t
}

func (d *Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

func (d *Date) Add(days int) Date {
	year, month, day := d.t.Date()
	return NewDate(year, int(month), day+days)
}

func (d *Date) Equal(ud Date) bool {
	return d.t.Equal(ud.Time())
}

// After reports whether d is after ud
func (d *Date) After(ud Date) bool {
	return d.t.After(ud.Time())
}

func (d *Date) AddDays(amount int) Date {
	year, month, date := d.t.Date()

	return NewDate(year, int(month), date+amount)
}

func ParseWeekday(wd string) (time.Weekday, bool) {
	switch lowerAndTrim(wd) {
	case "monday":
		return time.Monday, true
	case "tuesday":
		return time.Tuesday, true
	case "wednesday":
		return time.Wednesday, true
	case "thursday":
		return time.Thursday, true
	case "friday":
		return time.Friday, true
	case "saturday":
		return time.Saturday, true
	case "sunday":
		return time.Sunday, true
	case "maandag":
		return time.Monday, true
	case "dinsdag":
		return time.Tuesday, true
	case "woensdag":
		return time.Wednesday, true
	case "donderdag":
		return time.Thursday, true
	case "vrijdag":
		return time.Friday, true
	case "zaterdag":
		return time.Saturday, true
	case "zondag":
		return time.Sunday, true
	}

	return time.Monday, false
}

func lowerAndTrim(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}
