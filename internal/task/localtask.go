package task

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type LocalTask struct {
	Task
	LocalId     int
	LocalUpdate *LocalUpdate
}

func (lt *LocalTask) Apply(lu LocalUpdate) {
	if lu.Action != "" {
		lt.Action = lu.Action
	}
	if lu.Project != "" {
		lt.Project = lu.Project
	}
	if lu.Recur != nil {
		lt.Recur = lu.Recur
	}
	if !lu.Due.IsZero() {
		lt.Due = lu.Due
	}
	if lu.Done {
		lt.Done = lu.Done
	}
}

type ByDue []*LocalTask

func (lt ByDue) Len() int           { return len(lt) }
func (lt ByDue) Swap(i, j int)      { lt[i], lt[j] = lt[j], lt[i] }
func (lt ByDue) Less(i, j int) bool { return lt[j].Due.After(lt[i].Due) }

type ByDefault []*LocalTask

func (lt ByDefault) Len() int      { return len(lt) }
func (lt ByDefault) Swap(i, j int) { lt[i], lt[j] = lt[j], lt[i] }
func (lt ByDefault) Less(i, j int) bool {
	if !lt[j].Due.Equal(lt[i].Due) {
		return lt[j].Due.After(lt[i].Due)
	}

	if lt[i].Project != lt[j].Project {
		return lt[i].Project < lt[j].Project
	}

	return lt[i].LocalId < lt[j].LocalId
}

type LocalUpdate struct {
	ForVersion int
	Action     string
	Project    string
	Due        Date
	Recur      Recurrer
	Done       bool
}

func (lu LocalUpdate) Value() (driver.Value, error) {
	var recurStr string
	if lu.Recur != nil {
		recurStr = lu.Recur.String()
	}

	return fmt.Sprintf(`forversion: %d
action: %s
project: %s
recur: %s
due: %s
done: %t`,
		lu.ForVersion,
		lu.Action,
		lu.Project,
		recurStr,
		lu.Due.String(),
		lu.Done), nil
}

func (lu *LocalUpdate) Scan(value interface{}) error {
	body, err := driver.String.ConvertValue(value)
	if err != nil {
		*lu = LocalUpdate{}
		return nil
	}

	newLu := LocalUpdate{}
	for _, line := range strings.Split(body.(string), "\n") {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) < 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		switch k {
		case "forversion":
			d, _ := strconv.Atoi(v)
			newLu.ForVersion = d
		case "action":
			newLu.Action = v
		case "project":
			newLu.Project = v
		case "recur":
			newLu.Recur = NewRecurrer(v)
		case "due":
			newLu.Due = NewDateFromString(v)
		case "done":
			if v == "true" {
				newLu.Done = true
			}
		}
	}
	*lu = newLu

	return nil
}
