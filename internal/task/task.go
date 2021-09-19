package task

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ewintr.nl/gte/pkg/mstore"
	"github.com/google/uuid"
)

var (
	ErrOutdatedTask       = errors.New("task is outdated")
	ErrTaskIsNotRecurring = errors.New("task is not recurring")
)

const (
	FOLDER_INBOX     = "GTE/Inbox"
	FOLDER_NEW       = "GTE/New"
	FOLDER_RECURRING = "GTE/Recurring"
	FOLDER_PLANNED   = "GTE/Planned"
	FOLDER_UNPLANNED = "GTE/Unplanned"

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
	FIELD_DONE    = "done"
)

var (
	KnownFolders = []string{
		FOLDER_INBOX,
		FOLDER_NEW,
		FOLDER_RECURRING,
		FOLDER_PLANNED,
		FOLDER_UNPLANNED,
	}

	subjectFieldNames = []string{FIELD_ACTION}
	bodyFieldNames    = []string{
		FIELD_ID,
		FIELD_VERSION,
		FIELD_ACTION,
		FIELD_PROJECT,
		FIELD_DUE,
		FIELD_RECUR,
		FIELD_DONE,
	}
)

type Task struct {
	// Message is the underlying message from which the task was created
	// It only has meaning for remote repositories and will be nil in
	// local situations. It will be filtered out in LocalRepository.SetTasks()
	Message *mstore.Message

	Id      string
	Version int
	Folder  string

	Action  string
	Project string
	Due     Date
	Recur   Recurrer
	Done    bool
}

func NewFromMessage(msg *mstore.Message) *Task {
	t := &Task{
		Folder:  msg.Folder,
		Message: msg,
	}

	// parse fields from message
	subjectFields := map[string]string{}
	for _, f := range subjectFieldNames {
		subjectFields[f] = FieldFromSubject(f, msg.Subject)
	}

	bodyFields := map[string]string{}
	for _, f := range bodyFieldNames {
		value, _ := FieldFromBody(f, msg.Body)
		bodyFields[f] = value
	}

	// apply precedence rules
	version, _ := strconv.Atoi(bodyFields[FIELD_VERSION])
	id := bodyFields[FIELD_ID]
	if id == "" {
		id = uuid.New().String()
		version = 0
	}
	t.Id = id
	t.Version = version

	t.Action = bodyFields[FIELD_ACTION]
	if t.Action == "" {
		t.Action = subjectFields[FIELD_ACTION]
	}

	t.Project = bodyFields[FIELD_PROJECT]
	t.Due = NewDateFromString(bodyFields[FIELD_DUE])
	t.Recur = NewRecurrer(bodyFields[FIELD_RECUR])
	t.Done = bodyFields[FIELD_DONE] == "true"

	return t
}

func (t *Task) TargetFolder() string {
	switch {
	case t.Version == 0:
		return FOLDER_NEW
	case t.IsRecurrer():
		return FOLDER_RECURRING
	case !t.Due.IsZero():
		return FOLDER_PLANNED
	default:
		return FOLDER_UNPLANNED
	}
}

func (t *Task) NextMessage() *mstore.Message {
	tNew := t
	tNew.Folder = t.TargetFolder()
	tNew.Version++

	return &mstore.Message{
		Folder:  tNew.Folder,
		Subject: tNew.FormatSubject(),
		Body:    tNew.FormatBody(),
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
	if t.Done {
		fields[FIELD_DONE] = "true"
		order = append(order, FIELD_DONE)
	}

	keyLen := 0
	for _, f := range order {
		if len(f) > keyLen {
			keyLen = len(f)
		}
	}

	body := fmt.Sprintf("\r\n")
	for _, f := range order {
		key := f + FIELD_SEPARATOR
		for i := len(key); i <= keyLen; i++ {
			key += " "
		}
		line := strings.TrimSpace(fmt.Sprintf("%s %s", key, fields[f]))
		body += fmt.Sprintf("%s\r\n", line)
	}

	if t.Message != nil {
		body += fmt.Sprintf("\r\nPrevious version:\r\n\r\n%s\r\n", t.Message.Body)
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

func (t *Task) GenerateFromRecurrer(date Date) (*Task, error) {
	if !t.IsRecurrer() || !t.RecursOn(date) {
		return &Task{}, ErrTaskIsNotRecurring
	}

	return &Task{
		Id:      uuid.New().String(),
		Version: 1,
		Action:  t.Action,
		Project: t.Project,
		Due:     date,
	}, nil
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
