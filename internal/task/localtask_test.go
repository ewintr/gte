package task_test

import (
	"testing"

	"go-mod.ewintr.nl/go-kit/test"
	"go-mod.ewintr.nl/gte/internal/task"
)

func TestLocalTaskApply(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input *task.LocalTask
		exp   *task.LocalTask
	}{
		{
			name: "empty",
			input: &task.LocalTask{
				Task: task.Task{
					Action:  "action",
					Project: "project",
					Due:     task.NewDate(2021, 8, 22),
				},
				LocalUpdate: &task.LocalUpdate{},
			},
			exp: &task.LocalTask{
				Task: task.Task{
					Action:  "action",
					Project: "project",
					Due:     task.NewDate(2021, 8, 22),
				},
				LocalUpdate: &task.LocalUpdate{},
			},
		},
		{
			name: "all",
			input: &task.LocalTask{
				Task: task.Task{
					Version: 3,
				},
				LocalUpdate: &task.LocalUpdate{
					ForVersion: 3,
					Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_DONE},
					Action:     "action",
					Project:    "project",
					Due:        task.NewDate(2021, 8, 22),
					Done:       true,
				},
			},
			exp: &task.LocalTask{
				Task: task.Task{
					Version: 3,
					Action:  "action",
					Project: "project",
					Due:     task.NewDate(2021, 8, 22),
					Done:    true,
				},
				LocalUpdate: &task.LocalUpdate{},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.ApplyUpdate()
			test.Equals(t, tc.exp, tc.input)
		})
	}
}

func TestLocalUpdateAdd(t *testing.T) {
	for _, tc := range []struct {
		name  string
		start *task.LocalUpdate
		add   *task.LocalUpdate
		exp   *task.LocalUpdate
	}{
		{
			name:  "empty",
			start: &task.LocalUpdate{},
			add:   &task.LocalUpdate{},
			exp:   &task.LocalUpdate{},
		},
		{
			name: "empty add",
			start: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_RECUR, task.FIELD_DONE},
				Action:     "action",
				Project:    "project",
				Due:        task.NewDate(2021, 8, 22),
				Recur:      task.NewRecurrer("today, daily"),
				Done:       true,
			},
			add: &task.LocalUpdate{},
			exp: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_RECUR, task.FIELD_DONE},
				Action:     "action",
				Project:    "project",
				Due:        task.NewDate(2021, 8, 22),
				Recur:      task.NewRecurrer("today, daily"),
				Done:       true,
			},
		},
		{
			name:  "empty start",
			start: &task.LocalUpdate{},
			add: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_RECUR, task.FIELD_DONE},
				Action:     "action",
				Project:    "project",
				Due:        task.NewDate(2021, 8, 22),
				Recur:      task.NewRecurrer("today, daily"),
				Done:       true,
			},
			exp: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_RECUR, task.FIELD_DONE},
				Action:     "action",
				Project:    "project",
				Due:        task.NewDate(2021, 8, 22),
				Recur:      task.NewRecurrer("today, daily"),
				Done:       true,
			},
		},
		{
			name: "too old",
			start: &task.LocalUpdate{
				ForVersion: 3,
				Fields:     []string{},
			},
			add: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION},
			},
			exp: &task.LocalUpdate{
				ForVersion: 3,
				Fields:     []string{},
			},
		},
		{
			name: "adding fields",
			start: &task.LocalUpdate{
				ForVersion: 3,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT},
				Action:     "action-1",
				Project:    "project-1",
			},
			add: &task.LocalUpdate{
				ForVersion: 3,
				Fields:     []string{task.FIELD_PROJECT, task.FIELD_DUE},
				Project:    "project-2",
				Due:        task.NewDate(2021, 8, 22),
			},
			exp: &task.LocalUpdate{
				ForVersion: 3,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE},
				Action:     "action-1",
				Project:    "project-2",
				Due:        task.NewDate(2021, 8, 22),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.start.Add(tc.add)
			test.Equals(t, tc.exp, tc.start)
		})
	}
}
