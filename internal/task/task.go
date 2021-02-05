package task

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
	"github.com/google/uuid"
)

var (
	ErrOutdatedTask       = errors.New("task is outdated")
	ErrTaskIsNotRecurring = errors.New("task is not recurring")
)

const (
	FOLDER_INBOX     = "INBOX"
	FOLDER_NEW       = "New"
	FOLDER_RECURRING = "Recurring"
	FOLDER_PLANNED   = "Planned"
	FOLDER_UNPLANNED = "Unplanned"

	QUOTE_PREFIX       = ">"
	PREVIOUS_SEPARATOR = "Previous version:"

	FIELD_SEPARATOR   = ":"
	SUBJECT_SEPARATOR = " - "

	FIELD_ID      = "id"
	FIELD_VERSION = "version"
	FIELD_ACTION  = "action"
	FIELD_PROJECT = "project"
	FIELD_DUE     = "due"
	FIELD_RECUR   = "recur"
)

var (
	knownFolders = []string{
		FOLDER_INBOX,
		FOLDER_NEW,
		FOLDER_RECURRING,
		FOLDER_PLANNED,
		FOLDER_UNPLANNED,
	}
)

// Task reperesents a task based on the data stored in a message
type Task struct {

	// Id is a UUID that gets carried over when a new message is constructed
	Id string
	// Version is a method to determine the latest version for cleanup
	Version int

	// Folder is the same name as the mstore folder
	Folder string

	// Ordinary task attributes
	Action  string
	Project string
	Due     Date
	Recur   Recurrer

	//Message is the underlying message
	Message *mstore.Message

	// Current indicates whether the task represents an existing message in the mstore
	Current bool

	// Dirty indicates whether the task contains updates not present in the message
	Dirty bool
}

// New constructs a Task based on an mstore.Message.
//
// The data in the message is stored as key: value pairs, one per line. The line can start with quoting marks.
// The subject line also contains values in the format "date - project - action".
// Keys that exist more than once are merged. The one that appears first in the body takes precedence. A value present in the Body takes precedence over one in the subject.
// This enables updating a task by forwarding a topposted message whith new values for fields that the user wants to update.
func New(msg *mstore.Message) *Task {
	// Id
	dirty := false
	newId := false
	id, d := FieldFromBody(FIELD_ID, msg.Body)
	if id == "" {
		id = uuid.New().String()
		dirty = true
		newId = true
	}
	if d {
		dirty = true
	}

	// Version, cannot manually be incremented from body
	versionStr, _ := FieldFromBody(FIELD_VERSION, msg.Body)
	version, _ := strconv.Atoi(versionStr)
	if version == 0 {
		dirty = true
	}

	// Action
	action, d := FieldFromBody(FIELD_ACTION, msg.Body)
	if action == "" {
		action = FieldFromSubject(FIELD_ACTION, msg.Subject)
		if action != "" {
			dirty = true
		}
	}
	if d {
		dirty = true
	}

	// Due
	dueStr, d := FieldFromBody(FIELD_DUE, msg.Body)
	if dueStr == "" {
		dueStr = FieldFromSubject(FIELD_DUE, msg.Subject)
		if dueStr != "" {
			dirty = true
		}
	}
	if d {
		dirty = true
	}
	due := NewDateFromString(dueStr)

	// Recurrer
	recurStr, d := FieldFromBody(FIELD_RECUR, msg.Body)
	if d {
		dirty = true
	}
	recur := NewRecurrer(recurStr)

	// Folder
	folderOld := msg.Folder
	folderNew := folderOld
	if folderOld == FOLDER_INBOX {
		switch {
		case newId:
			folderNew = FOLDER_NEW
		case !newId && recur != nil:
			folderNew = FOLDER_RECURRING
		case !newId && recur == nil && due.IsZero():
			folderNew = FOLDER_UNPLANNED
		case !newId && recur == nil && !due.IsZero():
			folderNew = FOLDER_PLANNED
		}

	}
	if folderOld != folderNew {
		dirty = true
	}

	// Project
	project, d := FieldFromBody(FIELD_PROJECT, msg.Body)
	if project == "" {
		project = FieldFromSubject(FIELD_PROJECT, msg.Subject)
		if project != "" {
			dirty = true
		}
	}
	if d {
		dirty = true
	}

	if dirty {
		version++
	}

	return &Task{
		Id:      id,
		Version: version,
		Folder:  folderNew,
		Action:  action,
		Due:     due,
		Recur:   recur,
		Project: project,
		Message: msg,
		Current: true,
		Dirty:   dirty,
	}
}

func (t *Task) FormatSubject() string {
	var order []string
	if !t.Due.IsZero() {
		order = append(order, FIELD_DUE)
	}
	order = append(order, FIELD_PROJECT, FIELD_ACTION)

	fields := map[string]string{
		FIELD_PROJECT: t.Project,
		FIELD_ACTION:  t.Action,
		FIELD_DUE:     t.Due.String(),
	}

	parts := []string{}
	for _, f := range order {
		if fields[f] != "" {
			parts = append(parts, fields[f])
		}
	}

	return strings.Join(parts, SUBJECT_SEPARATOR)
}

func (t *Task) FormatBody() string {
	order := []string{FIELD_ACTION}
	fields := map[string]string{
		FIELD_ID:      t.Id,
		FIELD_VERSION: strconv.Itoa(t.Version),
		FIELD_PROJECT: t.Project,
		FIELD_ACTION:  t.Action,
	}
	if t.IsRecurrer() {
		order = append(order, FIELD_RECUR)
		fields[FIELD_RECUR] = t.Recur.String()
	} else {
		order = append(order, FIELD_DUE)
		fields[FIELD_DUE] = t.Due.String()
	}
	order = append(order, []string{FIELD_PROJECT, FIELD_VERSION, FIELD_ID}...)

	keyLen := 0
	for _, f := range order {
		if len(f) > keyLen {
			keyLen = len(f)
		}
	}

	body := fmt.Sprintf("\n")
	for _, f := range order {
		key := f + FIELD_SEPARATOR
		for i := len(key); i <= keyLen; i++ {
			key += " "
		}
		line := strings.TrimSpace(fmt.Sprintf("%s %s", key, fields[f]))
		body += fmt.Sprintf("%s\n", line)
	}

	if t.Message != nil {
		body += fmt.Sprintf("\nPrevious version:\n\n%s\n", t.Message.Body)
	}
	return body
}

func (t *Task) IsRecurrer() bool {
	return t.Recur != nil
}

func (t *Task) RecursToday() bool {
	return t.RecursOn(Today)
}

func (t *Task) RecursOn(date Date) bool {
	if !t.IsRecurrer() {
		return false
	}

	return t.Recur.RecursOn(date)
}

func (t *Task) CreateDueMessage(date Date) (string, string, error) {
	if !t.IsRecurrer() {
		return "", "", ErrTaskIsNotRecurring
	}

	tempTask := &Task{
		Id:      uuid.New().String(),
		Version: 1,
		Action:  t.Action,
		Project: t.Project,
		Due:     date,
	}

	return tempTask.FormatSubject(), tempTask.FormatBody(), nil
}

func FieldFromBody(field, body string) (string, bool) {
	value := ""
	dirty := false

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimPrefix(line, QUOTE_PREFIX))

		if line == PREVIOUS_SEPARATOR {
			return value, dirty
		}

		parts := strings.SplitN(line, FIELD_SEPARATOR, 2)
		if len(parts) < 2 {
			continue
		}

		fieldName := strings.ToLower(strings.TrimSpace(parts[0]))
		if fieldName == field {
			if value == "" {
				value = lowerAndTrim(parts[1])
			} else {
				dirty = true
			}
		}
	}

	return value, dirty
}

func FieldFromSubject(field, subject string) string {

	if field != FIELD_ACTION {
		return ""
	}

	terms := strings.Split(subject, SUBJECT_SEPARATOR)

	return lowerAndTrim(terms[len(terms)-1])
}
