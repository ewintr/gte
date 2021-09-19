package process_test

import (
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
)

func TestUpdate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		updates *task.LocalUpdate
	}{
		{
			name: "done",
			updates: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_DONE},
				Done:       true,
			},
		},
		{
			name: "fields",
			updates: &task.LocalUpdate{
				ForVersion: 2,
				Fields:     []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE},
				Project:    "project2",
				Action:     "action2",
				Due:        task.NewDate(2021, 8, 1),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			task1 := &task.Task{
				Id:      "id-1",
				Version: 2,
				Project: "project1",
				Action:  "action1",
				Due:     task.NewDate(2021, 7, 29),
				Folder:  task.FOLDER_PLANNED,
			}
			local := storage.NewMemory()
			allTasks := []*task.Task{task1}

			test.OK(t, local.SetTasks(allTasks))
			update := process.NewUpdate(local, task1.Id, tc.updates)
			test.OK(t, update.Process())
			lt, err := local.FindById(task1.Id)
			test.OK(t, err)
			test.Equals(t, task.STATUS_UPDATED, lt.LocalStatus)
			test.Equals(t, tc.updates, lt.LocalUpdate)
		})
	}
}
