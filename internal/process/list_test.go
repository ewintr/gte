package process_test

import (
	"errors"
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
)

func TestListProcess(t *testing.T) {
	date1 := task.NewDate(2021, 7, 9)
	date2 := task.NewDate(2021, 7, 10)
	date3 := task.NewDate(2021, 7, 11)

	task1 := &task.Task{
		Id:      "id1",
		Version: 1,
		Action:  "action1",
		Folder:  task.FOLDER_NEW,
	}
	task2 := &task.Task{
		Id:      "id2",
		Version: 1,
		Action:  "action2",
		Due:     date1,
		Folder:  task.FOLDER_PLANNED,
	}
	task3 := &task.Task{
		Id:      "id3",
		Version: 1,
		Action:  "action3",
		Due:     date2,
		Folder:  task.FOLDER_PLANNED,
	}
	task4 := &task.Task{
		Id:      "id4",
		Version: 1,
		Action:  "action4",
		Due:     date3,
		Folder:  task.FOLDER_PLANNED,
	}
	allTasks := []*task.Task{task1, task2, task3, task4}

	local := storage.NewMemory()
	test.OK(t, local.SetTasks(allTasks))

	t.Run("invalid reqs", func(t *testing.T) {
		list := process.NewList(local, process.ListReqs{})
		_, actErr := list.Process()
		test.Assert(t, errors.Is(actErr, process.ErrInvalidReqs), "expected invalid reqs err")
	})

	for _, tc := range []struct {
		name string
		reqs process.ListReqs
		exp  []*task.Task
	}{
		{
			name: "due",
			reqs: process.ListReqs{
				Due: date2,
			},
			exp: []*task.Task{task3},
		},
		{
			name: "due and before",
			reqs: process.ListReqs{
				Due:           date2,
				IncludeBefore: true,
			},
			exp: []*task.Task{task2, task3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			list := process.NewList(local, tc.reqs)

			act, err := list.Process()
			test.OK(t, err)
			test.Equals(t, tc.exp, act.Tasks)
		})
	}
}
