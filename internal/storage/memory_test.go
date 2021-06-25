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
		Folder:  folder1,
		Project: project1,
		Action:  "action1",
		Message: &mstore.Message{
			Subject: "action1",
		},
	}
	task2 := &task.Task{
		Folder:  folder1,
		Project: project2,
		Action:  "action2",
		Message: &mstore.Message{
			Subject: "action2",
		},
	}
	task3 := &task.Task{
		Folder:  folder2,
		Project: project1,
		Action:  "action3",
		Message: &mstore.Message{
			Subject: "action3",
		},
	}
	tasks := []*task.Task{task1, task2, task3}

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
		exp := []*task.Task{task1, task2}
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
		exp := []*task.Task{task1, task3}
		for _, tsk := range exp {
			tsk.Message = nil
		}
		test.Equals(t, exp, act)
	})
}
