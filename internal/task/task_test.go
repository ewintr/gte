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
	version := 2
	action := "some action"
	project := "project"
	folder := task.FOLDER_NEW

	for _, tc := range []struct {
		name       string
		message    *mstore.Message
		hasId      bool
		hasVersion bool
		exp        *task.Task
	}{
		{
			name:    "empty",
			message: &mstore.Message{},
			exp: &task.Task{
				Dirty: true,
			},
		},
		{
			name: "id, action, project and folder",
			message: &mstore.Message{
				Folder: folder,
				Body: fmt.Sprintf(`
id: %s
version: %d
action: %s
project: %s
`, id, version, action, project),
			},
			hasId:      true,
			hasVersion: true,
			exp: &task.Task{
				Id:      id,
				Version: version,
				Folder:  folder,
				Action:  action,
				Project: project,
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
version: %d
action: %s
				`, id, version, action),
			},
			hasId:      true,
			hasVersion: true,
			exp: &task.Task{
				Id:      id,
				Version: version,
				Folder:  folder,
				Action:  action,
			},
		},
		{
			name: "action from subject if not present in body",
			message: &mstore.Message{
				Folder:  folder,
				Subject: action,
				Body:    fmt.Sprintf(`id: %s`, id),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
				Dirty:  true,
			},
		},
		{
			name: "quoted fields",
			message: &mstore.Message{
				Folder: folder,
				Body: fmt.Sprintf(`
action: %s

Forwarded message:
> id: %s
> action: old action
				`, action, id),
			},
			hasId: true,
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
			if !tc.hasVersion {
				tc.exp.Version = 1
			}
			tc.exp.Message = tc.message
			tc.exp.Current = true
			test.Equals(t, tc.exp, act)
		})
	}
}

func TestFormatSubject(t *testing.T) {
	action := "an action"
	project := " a project"

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
			name: "action",
			task: &task.Task{Action: action},
			exp:  action,
		},
		{
			name: "project",
			task: &task.Task{Project: project},
			exp:  project,
		},
		{
			name: "action and project",
			task: &task.Task{Action: action, Project: project},
			exp:  fmt.Sprintf("%s - %s", project, action),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.task.FormatSubject())
		})
	}
}

func TestFormatBody(t *testing.T) {
	id := "an id"
	version := 6
	action := "an action"
	project := "project"

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
version: 0
project:
action:
`,
		},
		{
			name: "filled",
			task: &task.Task{
				Id:      id,
				Version: version,
				Action:  action,
				Project: project,
				Message: &mstore.Message{
					Body: "previous body",
				},
			},
			exp: `
id:      an id
version: 6
project: project
action:  an action

Previous version:

previous body
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
		name     string
		field    string
		body     string
		expValue string
		expDirty bool
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
			expValue: "valueb",
		},
		{
			name:  "present twice",
			field: "field",
			body: `
field: valuea
field: valueb
			`,
			expValue: "valuea",
			expDirty: true,
		},
		{
			name:     "colons",
			field:    "field",
			body:     "field:: val:ue",
			expValue: ": val:ue",
		},
		{
			name:  "trim field",
			field: "field",
			body: " field		: value",
			expValue: "value",
		},
		{
			name:  "trim value",
			field: "field",
			body: "field: 			value  ",
			expValue: "value",
		},
		{
			name: "quoted",

			field:    "field",
			body:     "> field: value",
			expValue: "value",
		},
		{
			name:  "previous body",
			field: "field",
			body: `
field: valuea

Previous version:

field: valueb
			`,
			expValue: "valuea",
		},
		{
			name:  "quoted previous body",
			field: "field",
			body: `
field: valuea

> Previous version:
>
> field: valueb
			`,
			expValue: "valuea",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actValue, actDirty := task.FieldFromBody(tc.field, tc.body)
			test.Equals(t, tc.expValue, actValue)
			test.Equals(t, tc.expDirty, actDirty)
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
