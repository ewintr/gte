package task_test

import (
	"fmt"
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func TestNewFromMessage(t *testing.T) {
	id := "an id"
	version := 2
	action := "some action"
	project := "project"
	date := task.NewDate(2021, 1, 20)

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
				Folder: task.FOLDER_UNPLANNED,
				Body: fmt.Sprintf(`
id: %s
due: no date
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
				Folder:  task.FOLDER_UNPLANNED,
				Action:  action,
				Project: project,
			},
		},
		{
			name: "with date",
			message: &mstore.Message{
				Folder: task.FOLDER_PLANNED,
				Body: fmt.Sprintf(`
id: %s
due: %s
version: %d
action: %s
`, id, date.String(), version, action),
			},
			hasId:      true,
			hasVersion: true,
			exp: &task.Task{
				Id:      id,
				Folder:  task.FOLDER_PLANNED,
				Action:  action,
				Version: version,
				Due:     date,
			},
		},
		{
			name: "folder inbox get updated to new",
			message: &mstore.Message{
				Folder: task.FOLDER_INBOX,
				Body: fmt.Sprintf(`
action: %s
`, action),
			},
			exp: &task.Task{
				Id:     id,
				Folder: task.FOLDER_NEW,
				Action: action,
				Dirty:  true,
			},
		},
		{
			name: "folder inbox gets updated to planned",
			message: &mstore.Message{
				Folder: task.FOLDER_INBOX,
				Body: fmt.Sprintf(`
id: %s
due: %s
action: %s
`, id, date.String(), action),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: task.FOLDER_PLANNED,
				Action: action,
				Due:    date,
				Dirty:  true,
			},
		},
		{
			name: "folder new gets updated to unplanned",
			message: &mstore.Message{
				Folder: task.FOLDER_INBOX,
				Body: fmt.Sprintf(`
id: %s
due: no date
action: %s
`, id, action),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: task.FOLDER_UNPLANNED,
				Action: action,
				Dirty:  true,
			},
		},
		{
			name: "action in body takes precedence",
			message: &mstore.Message{
				Folder:  task.FOLDER_PLANNED,
				Subject: "some other action",
				Body: fmt.Sprintf(`
id: %s
due: no date
version: %d
action: %s
				`, id, version, action),
			},
			hasId:      true,
			hasVersion: true,
			exp: &task.Task{
				Id:      id,
				Version: version,
				Folder:  task.FOLDER_PLANNED,
				Action:  action,
			},
		},
		{
			name: "action from subject if not present in body",
			message: &mstore.Message{
				Folder:  task.FOLDER_PLANNED,
				Subject: action,
				Body:    fmt.Sprintf(`id: %s`, id),
			},
			hasId: true,
			exp: &task.Task{
				Id:     id,
				Folder: task.FOLDER_PLANNED,
				Action: action,
				Dirty:  true,
			},
		},
		{
			name: "project in body takes precedence",
			message: &mstore.Message{
				Folder:  task.FOLDER_PLANNED,
				Subject: fmt.Sprintf("old project - %s", action),
				Body: fmt.Sprintf(`
id: %s
due: no date
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
				Folder:  task.FOLDER_PLANNED,
				Action:  action,
				Project: project,
			},
		},
		{
			name: "quoted fields",
			message: &mstore.Message{
				Folder: task.FOLDER_PLANNED,
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
				Folder: task.FOLDER_PLANNED,
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
	due := task.NewDate(2021, 1, 30)

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
			task: &task.Task{
				Action: action,
			},
			exp: action,
		},
		{
			name: "project",
			task: &task.Task{
				Project: project,
			},
			exp: project,
		},
		{
			name: "action and project",
			task: &task.Task{
				Action:  action,
				Project: project,
			},
			exp: fmt.Sprintf("%s - %s", project, action),
		},
		{
			name: "action, date and project",
			task: &task.Task{
				Action:  action,
				Project: project,
				Due:     due,
			},
			exp: fmt.Sprintf("%s - %s - %s", due.String(), project, action),
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
action:
due:     no date
project:
version: 0
id:
`,
		},
		{
			name: "filled",
			task: &task.Task{
				Id:      id,
				Version: version,
				Action:  action,
				Project: project,
				Due:     task.NewDate(2021, 1, 30),
				Message: &mstore.Message{
					Body: "previous body",
				},
			},
			exp: `
action:  an action
due:     2021-01-30 (saturday)
project: project
version: 6
id:      an id

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
	action := "action"
	project := "project"
	due := "due"
	subjectOne := fmt.Sprintf("%s", action)
	subjectTwo := fmt.Sprintf("%s - %s", project, action)
	subjectThree := fmt.Sprintf("%s - %s - %s", due, project, action)

	for _, tc := range []struct {
		name    string
		field   string
		subject string
		exp     string
	}{
		{
			name:    "empty field",
			subject: action,
		},
		{
			name:  "empty subject",
			field: task.FIELD_ACTION,
		},
		{
			name:    "unknown field",
			field:   "unknown",
			subject: subjectOne,
		},
		{
			name:    "action with one",
			field:   task.FIELD_ACTION,
			subject: subjectOne,
			exp:     action,
		},
		{
			name:    "action with with two",
			field:   task.FIELD_ACTION,
			subject: subjectTwo,
			exp:     action,
		},
		{
			name:    "action with three",
			field:   task.FIELD_ACTION,
			subject: subjectThree,
			exp:     action,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, task.FieldFromSubject(tc.field, tc.subject))
		})
	}
}
