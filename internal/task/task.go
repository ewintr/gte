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

	FIELD_SEPARATOR = ":"
	FIELD_ID        = "id"
	FIELD_ACTION    = "action"
)

// Task reperesents a task based on the data stored in a message
type Task struct {

	// Id is a UUID that gets carried over when a new message is constructed
	Id string

	// Folder is the same name as the mstore folder
	Folder string

	Action  string
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
	dirty := false
	id := FieldFromBody(FIELD_ID, msg.Body)
	if id == "" {
		id = uuid.New().String()
		dirty = true
	}

	action := FieldFromBody(FIELD_ACTION, msg.Body)
	if action == "" {
		action = FieldFromSubject(FIELD_ACTION, msg.Subject)
		if action != "" {
			dirty = true
		}
	}

	folder := msg.Folder
	if folder == FOLDER_INBOX {
		folder = FOLDER_NEW
		dirty = true
	}

	return &Task{
		Id:      id,
		Action:  action,
		Folder:  folder,
		Message: msg,
		Current: true,
		Dirty:   dirty,
	}
}

func (t *Task) FormatSubject() string {
	return t.Action
}

func (t *Task) FormatBody() string {
	body := fmt.Sprintf("\n")
	order := []string{FIELD_ID, FIELD_ACTION}
	fields := map[string]string{
		FIELD_ID:     t.Id,
		FIELD_ACTION: t.Action,
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

	return body
}

func FieldFromBody(field, body string) string {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, FIELD_SEPARATOR, 2)
		if len(parts) < 2 {
			continue
		}
		if strings.ToLower(strings.TrimSpace(parts[0])) == field {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

func FieldFromSubject(field, subject string) string {
	if field == FIELD_ACTION {
		return strings.ToLower(subject)
	}

	return ""
}
