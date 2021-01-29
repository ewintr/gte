package task_test

import (
	"fmt"
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

func TestNewFromMessage(t *testing.T) {
	id := "an id"
	action := "some action"
	folder := task.FOLDER_NEW
	for _, tc := range []struct {
		name    string
		message *mstore.Message
		hasId   bool
		exp     *task.Task
	}{
		{
			name:    "empty",
			message: &mstore.Message{},
			exp: &task.Task{
				Dirty: true,
			},
		},
		{
			name: "with id, action and folder",
			message: &mstore.Message{
				Folder: folder,
				Body: fmt.Sprintf(`
id: %s
action: %s
`, id, action),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
			},
		},
		{
			name: "folder inbox get updated to new",
			message: &mstore.Message{
				Folder: task.FOLDER_INBOX,
				Body: fmt.Sprintf(`
id: %s
action: %s
`, id, action),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: task.FOLDER_NEW,
				Action: action,
				Dirty:  true,
			},
		},
		{
			name: "action in subject takes precedence",
			message: &mstore.Message{
				Folder:  folder,
				Subject: "some other action",
				Body: fmt.Sprintf(`
id: %s
action: %s
				`, id, action),
			},
			exp: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
			},
		},
		{
			name: "action from subject if not present in body",
			message: &mstore.Message{
				Folder:  folder,
				Subject: action,
				Body:    fmt.Sprintf(`id: %s`, id),
			},
			exp: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
				Dirty:  true,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := task.New(tc.message)
			if !tc.hasId {
				test.Equals(t, false, "" == act.Id)
				tc.exp.Id = act.Id
			}
			tc.exp.Message = tc.message
			tc.exp.Current = true
			test.Equals(t, tc.exp, act)
		})
	}
}

func TestFormatSubject(t *testing.T) {
	action := "an action"
	for _, tc := range []struct {
		name string
		task *task.Task
		exp  string
	}{
		{
			name: "empty",
			task: &task.Task{},
		},
		{
			name: "with action",
			task: &task.Task{Action: action},
			exp:  action,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.task.FormatSubject())
		})
	}
}

func TestFormatBody(t *testing.T) {
	id := "an id"
	action := "an action"
	for _, tc := range []struct {
		name string
		task *task.Task
		exp  string
	}{
		{
			name: "empty",
			task: &task.Task{},
			exp: `
id:
action:
`,
		},
		{
			name: "filled",
			task: &task.Task{
				Id:     id,
				Action: action,
			},
			exp: `
id:     an id
action: an action
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.task.FormatBody())
		})
	}
}

func TestFieldFromBody(t *testing.T) {
	for _, tc := range []struct {
		name  string
		field string
		body  string
		exp   string
	}{
		{
			name: "empty field",
			body: `field: value`,
		},
		{
			name:  "empty body",
			field: "field",
		},
		{
			name:  "not present",
			field: "field",
			body:  "another: value",
		},
		{
			name:  "present",
			field: "fieldb",
			body: `
not a field at all

fielda: valuea
fieldb: valueb
fieldc: valuec
			`,
			exp: "valueb",
		},
		{
			name:  "present twice",
			field: "field",
			body: `
field: valuea
field: valueb
			`,
			exp: "valuea",
		},
		{
			name:  "with colons",
			field: "field",
			body:  "field:: val:ue",
			exp:   ": val:ue",
		},
		{
			name:  "trim field",
			field: "field",
			body: " field		: value",
			exp: "value",
		},
		{
			name:  "trim value",
			field: "field",
			body: "field: 			value  ",
			exp: "value",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, task.FieldFromBody(tc.field, tc.body))
		})
	}
}

func TestFieldFromSubject(t *testing.T) {
	for _, tc := range []struct {
		name    string
		field   string
		subject string
		exp     string
	}{
		{
			name:    "empty field",
			subject: "subject",
		},
		{
			name:  "empty subject",
			field: task.FIELD_ACTION,
		},
		{
			name:    "unknown field",
			field:   "unknown",
			subject: "subject",
		},
		{
			name:    "known field",
			field:   task.FIELD_ACTION,
			subject: "subject",
			exp:     "subject",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, task.FieldFromSubject(tc.field, tc.subject))
		})
	}
}
