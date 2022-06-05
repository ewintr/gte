package task

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

const (
	STATUS_FETCHED    = "fetched"
	STATUS_UPDATED    = "updated"
	STATUS_DISPATCHED = "dispatched"
)

type LocalTask struct {
	Task
	LocalId     int
	LocalUpdate *LocalUpdate
	LocalStatus string
}

func (lt *LocalTask) HasUpdate() bool {
	return lt.LocalUpdate.ForVersion != 0
}

func (lt *LocalTask) AddUpdate(update *LocalUpdate) {
	if lt.LocalUpdate == nil {
		lt.LocalUpdate = &LocalUpdate{}
	}

	lt.LocalUpdate.Add(update)
}

func (lt *LocalTask) ApplyUpdate() {
	if lt.LocalUpdate == nil {
		return
	}
	u := lt.LocalUpdate
	if u.ForVersion != lt.Version {
		lt.LocalUpdate = &LocalUpdate{}
		return
	}

	for _, field := range u.Fields {
		switch field {
		case FIELD_ACTION:
			lt.Action = u.Action
		case FIELD_PROJECT:
			lt.Project = u.Project
		case FIELD_DUE:
			lt.Due = u.Due
		case FIELD_RECUR:
			lt.Recur = u.Recur
		case FIELD_DONE:
			lt.Done = u.Done
		}
	}

	lt.LocalUpdate = &LocalUpdate{}
}

type ById []*LocalTask

func (lt ById) Len() int           { return len(lt) }
func (lt ById) Swap(i, j int)      { lt[i], lt[j] = lt[j], lt[i] }
func (lt ById) Less(i, j int) bool { return lt[i].Id < lt[j].Id }

type ByDue []*LocalTask

func (lt ByDue) Len() int           { return len(lt) }
func (lt ByDue) Swap(i, j int)      { lt[i], lt[j] = lt[j], lt[i] }
func (lt ByDue) Less(i, j int) bool { return lt[j].Due.After(lt[i].Due) }

type ByDefault []*LocalTask

func (lt ByDefault) Len() int      { return len(lt) }
func (lt ByDefault) Swap(i, j int) { lt[i], lt[j] = lt[j], lt[i] }
func (lt ByDefault) Less(i, j int) bool {
	if lt[i].IsRecurrer() != lt[j].IsRecurrer() {
		return lt[i].IsRecurrer()
	}
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
	Fields     []string
	Action     string
	Project    string
	Due        Date
	Recur      Recurrer
	Done       bool
}

func (lu *LocalUpdate) Add(newUpdate *LocalUpdate) {
	if lu.ForVersion > newUpdate.ForVersion {
		return
	}
	lu.ForVersion = newUpdate.ForVersion

	for _, nf := range newUpdate.Fields {
		switch nf {
		case FIELD_ACTION:
			lu.Action = newUpdate.Action
		case FIELD_PROJECT:
			lu.Project = newUpdate.Project
		case FIELD_DUE:
			lu.Due = newUpdate.Due
		case FIELD_RECUR:
			lu.Recur = newUpdate.Recur
		case FIELD_DONE:
			lu.Done = newUpdate.Done
		}

		add := true
		for _, of := range lu.Fields {
			if nf == of {
				add = false
				break
			}
		}
		if add {
			lu.Fields = append(lu.Fields, nf)
		}
	}
}

func (lu LocalUpdate) Value() (driver.Value, error) {
	v := fmt.Sprintf("forversion: %d\n", lu.ForVersion)
	for _, f := range lu.Fields {
		switch f {
		case FIELD_ACTION:
			v += fmt.Sprintf("action: %s\n", lu.Action)
		case FIELD_PROJECT:
			v += fmt.Sprintf("project: %s\n", lu.Project)
		case FIELD_RECUR:
			v += fmt.Sprintf("recur: %s\n", lu.Recur.String())
		case FIELD_DUE:
			v += fmt.Sprintf("due: %s\n", lu.Due)
		case FIELD_DONE:
			v += fmt.Sprintf("done: %t\n", lu.Done)
		}
	}

	return v, nil
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
			newLu.Fields = append(newLu.Fields, FIELD_ACTION)
		case "project":
			newLu.Project = v
			newLu.Fields = append(newLu.Fields, FIELD_PROJECT)
		case "recur":
			newLu.Recur = NewRecurrer(v)
			newLu.Fields = append(newLu.Fields, FIELD_RECUR)
		case "due":
			newLu.Due = NewDateFromString(v)
			newLu.Fields = append(newLu.Fields, FIELD_DUE)
		case "done":
			if v == "true" {
				newLu.Done = true
				newLu.Fields = append(newLu.Fields, FIELD_DONE)
			}
		}
	}
	*lu = newLu

	return nil
}
