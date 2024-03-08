package process_test

import (
	"testing"

	"code.ewintr.nl/go-kit/test"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
)

func TestNew(t *testing.T) {
	local := storage.NewMemory()
	update := &task.LocalUpdate{
		Fields:  []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE},
		Project: "project",
		Action:  "action",
		Due:     task.NewDate(2021, 9, 4),
	}
	n := process.NewNew(local, update)
	test.OK(t, n.Process())
	tasks, err := local.FindAll()
	test.OK(t, err)
	test.Assert(t, len(tasks) == 1, "amount of tasks was not 1")
	tsk := tasks[0]
	test.Assert(t, tsk.Id != "", "id was empty")
	test.Equals(t, update, tsk.LocalUpdate)
}
