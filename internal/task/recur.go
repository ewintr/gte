package task

import (
	"fmt"
	"strings"
	"time"
)

type Recurrer interface {
	RecursOn(date Date) bool
	String() string
}

func NewRecurrer(recurStr string) Recurrer {
	terms := strings.Split(recurStr, ", ")
	if len(terms) < 2 {
		return nil
	}

	start := NewDateFromString(terms[0])
	if start.IsZero() {
		return nil
	}

	terms = terms[1:]

	if recur, ok := ParseDaily(start, terms); ok {
		return recur
	}
	if recur, ok := ParseWeekly(start, terms); ok {
		return recur
	}
	if recur, ok := ParseBiweekly(start, terms); ok {
		return recur
	}

	return nil
}

type Daily struct {
	Start Date
}

func ParseDaily(start Date, terms []string) (Recurrer, bool) {
	if len(terms) < 1 {
		return nil, false
	}

	if terms[0] != "daily" {
		return nil, false
	}

	return Daily{
		Start: start,
	}, true
}

func (d Daily) RecursOn(date Date) bool {
	return date.Equal(d.Start) || date.After(d.Start)
}

func (d Daily) String() string {
	return fmt.Sprintf("%s, daily", d.Start.String())
}

type Weekly struct {
	Start    Date
	Weekdays Weekdays
}

// yyyy-mm-dd, weekly, wednesday & saturday & sunday
func ParseWeekly(start Date, terms []string) (Recurrer, bool) {
	if len(terms) < 2 {
		return nil, false
	}

	if terms[0] != "weekly" {
		return nil, false
	}

	wds := Weekdays{}
	for _, wdStr := range strings.Split(terms[1], "&") {
		wd, ok := ParseWeekday(wdStr)
		if !ok {
			continue
		}
		wds = append(wds, wd)
	}
	if len(wds) == 0 {
		return nil, false
	}

	return Weekly{
		Start:    start,
		Weekdays: wds.Unique(),
	}, true
}

func (w Weekly) RecursOn(date Date) bool {
	if w.Start.After(date) {
		return false
	}

	for _, wd := range w.Weekdays {
		if wd == date.Weekday() {
			return true
		}
	}

	return false
}

func (w Weekly) String() string {
	weekdayStrs := []string{}
	for _, wd := range w.Weekdays {
		weekdayStrs = append(weekdayStrs, wd.String())
	}
	weekdayStr := strings.Join(weekdayStrs, " & ")

	return fmt.Sprintf("%s, weekly, %s", w.Start.String(), strings.ToLower(weekdayStr))
}

type Biweekly struct {
	Start   Date
	Weekday time.Weekday
}

// yyyy-mm-dd, biweekly, wednesday
func ParseBiweekly(start Date, terms []string) (Recurrer, bool) {
	if len(terms) < 2 {
		return nil, false
	}

	if terms[0] != "biweekly" {
		return nil, false
	}

	wd, ok := ParseWeekday(terms[1])
	if !ok {
		return nil, false
	}

	return Biweekly{
		Start:   start,
		Weekday: wd,
	}, true
}

func (b Biweekly) RecursOn(date Date) bool {
	if b.Start.After(date) {
		return false
	}

	if b.Weekday != date.Weekday() {
		return false
	}

	// find first
	tDate := b.Start
	for {
		if tDate.Weekday() == b.Weekday {
			break
		}
		tDate = tDate.AddDays(1)
	}

	// add weeks
	for {
		switch {
		case tDate.Equal(date):
			return true
		case tDate.After(date):
			return false
		}
		tDate = tDate.AddDays(14)
	}
}

func (b Biweekly) String() string {
	return fmt.Sprintf("%s, biweekly, %s", b.Start.String(), strings.ToLower(b.Weekday.String()))
}