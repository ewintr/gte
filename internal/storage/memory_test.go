package storage_test

import (
	"testing"
	"time"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func TestMemory(t *testing.T) {
	folder1, folder2 := "folder1", "folder2"
	project1, project2 := "project1", "project2"
	task1 := &task.Task{
		Id:      "id-1",
		Folder:  folder1,
		Project: project1,
		Action:  "action1",
		Message: &mstore.Message{
			Subject: "action1",
		},
	}
	task2 := &task.Task{
		Id:      "id-2",
		Folder:  folder1,
		Project: project2,
		Action:  "action2",
		Message: &mstore.Message{
			Subject: "action2",
		},
	}
	task3 := &task.Task{
		Id:      "id-3",
		Folder:  folder2,
		Project: project1,
		Action:  "action3",
		Message: &mstore.Message{
			Subject: "action3",
		},
	}
	tasks := []*task.Task{task1, task2, task3}
	localTask1 := &task.LocalTask{Task: *task1, LocalId: 1}
	localTask2 := &task.LocalTask{Task: *task2, LocalId: 2}
	localTask3 := &task.LocalTask{Task: *task3, LocalId: 3}

	t.Run("sync", func(t *testing.T) {
		mem := storage.NewMemory()
		latest, err := mem.LatestSync()
		test.OK(t, err)
		test.Assert(t, latest.IsZero(), "lastest was not zero")

		start := time.Now()
		test.OK(t, mem.SetTasks(tasks))
		latest, err = mem.LatestSync()
		test.OK(t, err)
		test.Assert(t, latest.After(start), "latest was not after start")
	})

	t.Run("findallinfolder", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindAllInFolder(folder1)
		test.OK(t, err)
		exp := []*task.LocalTask{localTask1, localTask2}
		for _, tsk := range exp {
			tsk.Message = nil
		}
		test.Equals(t, exp, act)
	})

	t.Run("findallinproject", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindAllInProject(project1)
		test.OK(t, err)
		exp := []*task.LocalTask{localTask1, localTask3}
		for _, tsk := range exp {
			tsk.Message = nil
		}
		test.Equals(t, exp, act)
	})

	t.Run("findbyid", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindById("id-2")
		test.OK(t, err)
		test.Equals(t, localTask2, act)
	})

	t.Run("findbylocalid", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindByLocalId(2)
		test.OK(t, err)
		test.Equals(t, localTask2, act)
	})

	t.Run("setlocalupdate", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		expUpdate := &task.LocalUpdate{
			ForVersion: 1,
			Action:     "update action",
			Project:    "update project",
			Due:        task.NewDate(2021, 8, 21),
			Recur:      task.NewRecurrer("today, weekly, monday"),
			Done:       true,
		}
		test.OK(t, mem.SetLocalUpdate(2, expUpdate))
		actTask, err := mem.FindByLocalId(2)
		test.OK(t, err)
		test.Equals(t, expUpdate, actTask.LocalUpdate)
	})
}
