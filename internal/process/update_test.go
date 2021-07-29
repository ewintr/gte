package process_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

func TestUpdate(t *testing.T) {
	task1 := &task.Task{
		Id:      "id-1",
		Project: "project1",
		Action:  "action1",
		Due:     task.NewDate(2021, 7, 29),
		Folder:  task.FOLDER_PLANNED,
	}
	local := storage.NewMemory()
	allTasks := []*task.Task{task1}

	for _, tc := range []struct {
		name    string
		updates process.UpdateFields
		exp     *task.Task
	}{
		{
			name: "done",
			updates: process.UpdateFields{
				task.FIELD_DONE: "true",
			},
			exp: &task.Task{
				Id:      "id-1",
				Project: "project1",
				Action:  "action1",
				Due:     task.NewDate(2021, 7, 29),
				Folder:  task.FOLDER_PLANNED,
				Done:    true,
			},
		},
		{
			name: "fields",
			updates: process.UpdateFields{
				task.FIELD_PROJECT: "project2",
				task.FIELD_ACTION:  "action2",
				task.FIELD_DUE:     "2021-08-01",
			},
			exp: &task.Task{
				Id:      "id-1",
				Project: "project2",
				Action:  "action2",
				Due:     task.NewDate(2021, 8, 1),
				Folder:  task.FOLDER_PLANNED,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			local.SetTasks(allTasks)
			out := msend.NewMemory()
			disp := storage.NewDispatcher(out)

			update := process.NewUpdate(local, disp, task1.Id, tc.updates)
			test.OK(t, update.Process())
			expMsg := &msend.Message{
				Subject: tc.exp.FormatSubject(),
				Body:    tc.exp.FormatBody(),
			}
			test.Assert(t, len(out.Messages) == 1, "amount of messages was not one")
			test.Equals(t, expMsg, out.Messages[0])
		})
	}
}
