package storage_test

import (
	"fmt"
	"testing"

	"go-mod.ewintr.nl/go-kit/test"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
	"go-mod.ewintr.nl/gte/pkg/msend"
)

func TestDispatcherDispatch(t *testing.T) {
	mem := msend.NewMemory()
	disp := storage.NewDispatcher(mem)
	tsk := &task.Task{
		Id:      "id",
		Version: 3,
		Action:  "action",
		Project: "project",
		Due:     task.NewDate(2021, 6, 24),
	}

	t.Run("err", func(t *testing.T) {
		expErr := fmt.Errorf("not good")
		mem.Err = expErr
		actErr := disp.Dispatch(tsk)

		test.Equals(t, expErr, actErr)
	})

	t.Run("success", func(t *testing.T) {
		mem.Err = nil

		test.OK(t, disp.Dispatch(tsk))
		test.Equals(t, 1, len(mem.Messages))

		actMsg := mem.Messages[0]
		test.Equals(t, tsk.FormatSubject(), actMsg.Subject)
		test.Equals(t, tsk.FormatBody(), actMsg.Body)
	})
}
