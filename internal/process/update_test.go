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
		Folder:  task.FOLDER_PLANNED,
	}
	local := storage.NewMemory()
	out := msend.NewMemory()
	disp := storage.NewDispatcher(out)
	allTasks := []*task.Task{task1}

	t.Run("done", func(t *testing.T) {
		local.SetTasks(allTasks)
		updates := process.UpdateFields{
			"done": "true",
		}

		update := process.NewUpdate(local, disp, task1.Id, updates)
		test.OK(t, update.Process())
		expTask := task1
		expTask.Done = true
		expMsg := &msend.Message{
			Subject: expTask.FormatSubject(),
			Body:    expTask.FormatBody(),
		}
		test.Assert(t, len(out.Messages) == 1, "amount of messages was not one")
		test.Equals(t, expMsg, out.Messages[0])

	})
}
