package task_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/task"
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
