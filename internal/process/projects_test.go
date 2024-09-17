package process_test

import (
	"testing"

	"go-mod.ewintr.nl/go-kit/test"
	"go-mod.ewintr.nl/gte/internal/process"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
)

func TestProjects(t *testing.T) {
	project1, project2, project3 := "project-1", "project-2", "project-3"

	task1 := &task.Task{
		Id:      "id1",
		Version: 1,
		Action:  "action1",
		Folder:  task.FOLDER_NEW,
		Project: project1,
	}
	task2 := &task.Task{
		Id:      "id2",
		Version: 1,
		Action:  "action2",
		Due:     task.NewDate(2021, 8, 19),
		Folder:  task.FOLDER_PLANNED,
		Project: project2,
	}
	task3 := &task.Task{
		Id:      "id3",
		Version: 1,
		Action:  "action3",
		Due:     task.NewDate(2021, 8, 18),
		Folder:  task.FOLDER_PLANNED,
		Project: project2,
	}
	task4 := &task.Task{
		Id:      "id4",
		Version: 1,
		Action:  "action4",
		Due:     task.NewDate(2021, 8, 17),
		Folder:  task.FOLDER_UNPLANNED,
		Project: project3,
	}
	allTasks := []*task.Task{task1, task2, task3, task4}

	local := storage.NewMemory()
	test.OK(t, local.SetTasks(allTasks))

	t.Run("all", func(t *testing.T) {
		exp := []string{project1, project2, project3}
		list := process.NewProjects(local)
		act, err := list.Process()
		test.OK(t, err)
		test.Equals(t, exp, act)
	})
}
