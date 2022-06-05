package task

import (
	"fmt"
	"sort"
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

type Weekdays []time.Weekday

func (wds Weekdays) Len() int      { return len(wds) }
func (wds Weekdays) Swap(i, j int) { wds[j], wds[i] = wds[i], wds[j] }
func (wds Weekdays) Less(i, j int) bool {
	if wds[i] == time.Sunday {
		return false
	}
	if wds[j] == time.Sunday {
		return true
	}

	return int(wds[i]) < int(wds[j])
}

func (wds Weekdays) Unique() Weekdays {
	mwds := map[time.Weekday]bool{}
	for _, wd := range wds {
		mwds[wd] = true
	}
	newWds := Weekdays{}
	for wd := range mwds {
		newWds = append(newWds, wd)
	}
	sort.Sort(newWds)

	return newWds
}

type Date struct {
	t time.Time
}

func NewDate(year, month, day int) Date {

	if year == 0 || month == 0 || month > 12 || day == 0 {
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

	switch date {
	case "":
		fallthrough
	case "no-date":
		fallthrough
	case "no date":
		return Date{}

	case "today":
		return Today

	case "tomorrow":
		return Today.AddDays(1)

	case "day after tomorrow":
		fallthrough
	case "day-after-tomorrow":
		return Today.AddDays(2)

	case "this week":
		fallthrough
	case "this-week":
		daysToAdd := findDaysToWeekday(Today.Weekday(), time.Friday)
		return Today.Add(daysToAdd)

	case "next week":
		fallthrough
	case "next-week":
		daysToAdd := findDaysToWeekday(Today.Weekday(), time.Friday) + 7
		return Today.Add(daysToAdd)

	case "this sprint":
		tDate := NewDate(2021, 1, 28) // a sprint end
		for {
			if tDate.After(Today) {
				return tDate
			}
			tDate = tDate.AddDays(14)
		}
	}

	t, err := time.Parse("2006-01-02", fmt.Sprintf("%.10s", date))
	if err == nil {
		return Date{t: t}
	}

	newWeekday, ok := ParseWeekday(date)
	if !ok {
		return Date{}
	}
	daysToAdd := findDaysToWeekday(Today.Weekday(), newWeekday)

	return Today.Add(daysToAdd)
}

func findDaysToWeekday(current, wanted time.Weekday) int {
	daysToAdd := int(wanted) - int(current)
	if daysToAdd <= 0 {
		daysToAdd += 7
	}

	return daysToAdd
}

func (d Date) DaysBetween(d2 Date) int {
	tDate := d2
	end := d
	if !end.After(tDate) {
		end = d2
		tDate = d
	}

	days := 0
	for {
		if tDate.Add(days).Equal(end) {
			return days
		}
		days++
	}
}

func (d Date) String() string {
	if d.t.IsZero() {
		return "no date"
	}

	return strings.ToLower(d.t.Format(DateFormat))
}

func (d Date) Human() string {
	if d.t.IsZero() {
		return "-"
	}

	if Today.Add(7).After(d) {
		return strings.ToLower(d.t.Format("Monday"))
	}

	return strings.ToLower(d.t.Format(DateFormat))
}

func (d Date) IsZero() bool {
	return d.t.IsZero()
}

func (d Date) Time() time.Time {
	return d.t
}

func (d Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

func (d Date) Day() int {
	return d.t.Day()
}

func (d Date) Add(days int) Date {
	year, month, day := d.t.Date()
	return NewDate(year, int(month), day+days)
}

func (d Date) AddMonths(addMonths int) Date {
	year, mmonth, day := d.t.Date()
	month := int(mmonth)
	for m := 1; m <= addMonths; m++ {
		month += 1
		if month == 12 {
			year += 1
			month = 1
		}
	}

	return NewDate(year, month, day)
}

func (d Date) Equal(ud Date) bool {
	return d.t.Equal(ud.Time())
}

// After reports whether d is after ud
func (d Date) After(ud Date) bool {
	return d.t.After(ud.Time())
}

func (d Date) AddDays(amount int) Date {
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
