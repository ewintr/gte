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
		Project: "project1",
	}
	task2 := &task.Task{
		Id:      "id2",
		Version: 1,
		Action:  "action2",
		Due:     date1,
		Folder:  task.FOLDER_PLANNED,
		Project: "project2",
	}
	task3 := &task.Task{
		Id:      "id3",
		Version: 1,
		Action:  "action3",
		Due:     date2,
		Folder:  task.FOLDER_PLANNED,
		Project: "project1",
	}
	task4 := &task.Task{
		Id:      "id4",
		Version: 1,
		Action:  "action4",
		Due:     date3,
		Folder:  task.FOLDER_PLANNED,
		Project: "project2",
	}
	allTasks := []*task.Task{task1, task2, task3, task4}
	localTask2 := &task.LocalTask{Task: *task2, LocalId: 2}
	localTask3 := &task.LocalTask{Task: *task3, LocalId: 3}
	localTask4 := &task.LocalTask{Task: *task4, LocalId: 4}
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
		exp  []*task.LocalTask
	}{
		{
			name: "due",
			reqs: process.ListReqs{
				Due: date2,
			},
			exp: []*task.LocalTask{localTask3},
		},
		{
			name: "due and before",
			reqs: process.ListReqs{
				Due:           date2,
				IncludeBefore: true,
			},
			exp: []*task.LocalTask{localTask2, localTask3},
		},
		{
			name: "folder",
			reqs: process.ListReqs{
				Folder: task.FOLDER_PLANNED,
			},
			exp: []*task.LocalTask{localTask2, localTask3, localTask4},
		},
		{
			name: "project",
			reqs: process.ListReqs{
				Project: "project2",
			},
			exp: []*task.LocalTask{localTask2, localTask4},
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
