package task

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type LocalTask struct {
	Task
	LocalId int
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

type LocalUpdate struct {
	Action  string
	Project string
	Due     Date
	Recur   Recurrer
	Done    bool
}

func (lu LocalUpdate) Value() (driver.Value, error) {
	return fmt.Sprintf(`action: %s
project: %s
recur: %s
due: %s
done: %t`,
		lu.Action,
		lu.Project,
		lu.Recur.String(),
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
