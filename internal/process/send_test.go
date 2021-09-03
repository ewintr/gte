package process_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
)

func TestSend(t *testing.T) {
	task1 := &task.Task{
		Id:      "id-1",
		Version: 2,
		Project: "project1",
		Action:  "action1",
		Due:     task.NewDate(2021, 7, 29),
		Folder:  task.FOLDER_PLANNED,
	}
	task2 := &task.Task{
		Id:      "id-2",
		Version: 2,
		Project: "project1",
		Action:  "action2",
		Folder:  task.FOLDER_UNPLANNED,
	}
	local := storage.NewMemory()
	allTasks := []*task.Task{task1, task2}

	test.OK(t, local.SetTasks(allTasks))

	t.Run("no updates", func(t *testing.T) {
		out := msend.NewMemory()
		disp := storage.NewDispatcher(out)
		send := process.NewSend(local, disp)
		test.OK(t, send.Process())
		test.Assert(t, len(out.Messages) == 0, "amount of messages was not 0")
	})

	t.Run("update", func(t *testing.T) {
		lu := &task.LocalUpdate{
			ForVersion: task2.Version,
			Fields:     []string{task.FIELD_ACTION},
			Action:     "updated",
		}
		lt, err := local.FindById(task2.Id)
		test.OK(t, err)
		lt.AddUpdate(lu)
		test.OK(t, local.SetLocalUpdate(lt))

		out := msend.NewMemory()
		disp := storage.NewDispatcher(out)
		send := process.NewSend(local, disp)
		test.OK(t, send.Process())
		test.Assert(t, len(out.Messages) == 1, "amount of messages was not 1")
		expSubject := "project1 - updated"
		test.Equals(t, expSubject, out.Messages[0].Subject)
	})
}
