package task

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
	"github.com/google/uuid"
)

var (
	ErrOutdatedTask = errors.New("task is outdated")
)

const (
	FOLDER_INBOX = "INBOX"
	FOLDER_NEW   = "New"

	QUOTE_PREFIX       = ">"
	PREVIOUS_SEPARATOR = "Previous version:"

	FIELD_SEPARATOR   = ":"
	SUBJECT_SEPARATOR = " - "

	FIELD_ID      = "id"
	FIELD_ACTION  = "action"
	FIELD_PROJECT = "project"
	FIELD_DUE     = "date"
)

var (
	knownFolders = []string{FOLDER_INBOX, FOLDER_NEW}
)

// Task reperesents a task based on the data stored in a message
type Task struct {

	// Id is a UUID that gets carried over when a new message is constructed
	Id string

	// Folder is the same name as the mstore folder
	Folder string

	// Ordinary task attributes
	Action  string
	Project string
	Due     Date

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
	id, d := FieldFromBody(FIELD_ID, msg.Body)
	if id == "" {
		id = uuid.New().String()
		dirty = true
	}
	if d {
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

	// Folder
	folder := msg.Folder
	if folder == FOLDER_INBOX {
		folder = FOLDER_NEW
		dirty = true
	}

	// Project
	project, d := FieldFromBody(FIELD_PROJECT, msg.Body)
	if d {
		dirty = true
	}

	return &Task{
		Id:      id,
		Folder:  folder,
		Action:  action,
		Project: project,
		Message: msg,
		Current: true,
		Dirty:   dirty,
	}
}

func (t *Task) FormatSubject() string {
	order := []string{FIELD_PROJECT, FIELD_ACTION}
	fields := map[string]string{
		FIELD_PROJECT: t.Project,
		FIELD_ACTION:  t.Action,
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
	body := fmt.Sprintf("\n")
	order := []string{FIELD_ID, FIELD_PROJECT, FIELD_ACTION}
	fields := map[string]string{
		FIELD_ID:      t.Id,
		FIELD_PROJECT: t.Project,
		FIELD_ACTION:  t.Action,
	}

	keyLen := 0
	for _, f := range order {
		if len(f) > keyLen {
			keyLen = len(f)
		}
	}

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
				value = strings.TrimSpace(parts[1])
			} else {
				dirty = true
			}
		}
	}

	return value, dirty
}

func FieldFromSubject(field, subject string) string {
	if field == FIELD_ACTION {
		return strings.ToLower(subject)
	}

	return ""
}
