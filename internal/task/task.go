package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
	"github.com/google/uuid"
)

var (
	ErrOutdatedTask = errors.New("task is outdated")
)

type Date time.Time

func (d *Date) Weekday() Weekday {
	return d.Weekday()
}

type Task struct {
	Id         string
	Folder     string
	Action     string
	Due        Date
	Message    *mstore.Message
	Current    bool
	Simplified bool
}

func NewFromMessage(msg *mstore.Message) *Task {
	fmt.Println(msg.Subject)
	id := FieldFromBody("id", msg.Body)
	if id == "" {
		id = uuid.New().String()
	}

	action := FieldFromBody("action", msg.Body)
	if action == "" {
		action = FieldFromSubject("action", msg.Subject)
	}

	folder := msg.Folder
	if folder == "INBOX" {
		folder = "New"
	}

	return &Task{
		Id:         id,
		Action:     action,
		Folder:     folder,
		Message:    msg,
		Current:    true,
		Simplified: false,
	}
}

// Dirty checks if the task has unsaved changes
func (t *Task) Dirty() bool {
	mBody := t.Message.Body
	mSubject := t.Message.Subject

	if t.Id != FieldFromBody("id", mBody) {
		return true
	}

	if t.Folder != t.Message.Folder {
		return true
	}

	if t.Action != FieldFromBody("action", mBody) {
		return true
	}
	if t.Action != FieldFromSubject("action", mSubject) {
		return true
	}

	return false
}

func (t *Task) Subject() string {
	return t.Action
}

func (t *Task) Body() string {
	body := fmt.Sprintf("id: %s\n", t.Id)
	body += fmt.Sprintf("action: %s\n", t.Action)

	return body
}

func FieldFromBody(field, body string) string {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		if strings.ToLower(parts[0]) == field {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

func FieldFromSubject(field, subject string) string {
	if field == "action" {
		return strings.ToLower(subject)
	}

	return ""
}
