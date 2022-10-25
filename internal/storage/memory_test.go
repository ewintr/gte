package storage_test

import (
	"sort"
	"testing"
	"time"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
	"ewintr.nl/gte/pkg/mstore"
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
	emptyUpdate := &task.LocalUpdate{}
	localTask1 := &task.LocalTask{Task: *task1, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED}
	localTask2 := &task.LocalTask{Task: *task2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED}
	localTask3 := &task.LocalTask{Task: *task3, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED}

	t.Run("sync", func(t *testing.T) {
		mem := storage.NewMemory(task1)
		latestFetch, latestDisp, err := mem.LatestSyncs()
		test.OK(t, err)
		test.Assert(t, latestFetch.IsZero(), "latestfetch was not zero")
		test.Assert(t, latestDisp.IsZero(), "latestdisp  was not zero")

		start := time.Now()
		test.OK(t, mem.SetTasks(tasks))
		test.OK(t, mem.MarkDispatched(1))
		latestFetch, latestDisp, err = mem.LatestSyncs()
		test.OK(t, err)
		test.Assert(t, latestFetch.After(start), "latestfetch was not after start")
		test.Assert(t, latestDisp.After(start), "latestdisp was not after start")
	})

	t.Run("findallin", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindAll()
		test.OK(t, err)
		exp := []*task.LocalTask{localTask1, localTask2, localTask3}
		for _, tsk := range act {
			tsk.LocalId = 0
		}
		sExp := task.ById(exp)
		sAct := task.ById(act)
		sort.Sort(sExp)
		sort.Sort(sAct)
		test.Equals(t, sExp, sAct)
	})

	t.Run("findbyid", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		act, err := mem.FindById("id-2")
		test.OK(t, err)
		act.LocalId = 0
		test.Equals(t, localTask2, act)
	})

	t.Run("findbylocalid", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks([]*task.Task{task1}))
		act, err := mem.FindByLocalId(1)
		test.OK(t, err)
		act.LocalId = 0
		test.Equals(t, localTask1, act)
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
		test.OK(t, mem.SetLocalUpdate(task2.Id, expUpdate))
		actTask, err := mem.FindById(task2.Id)
		test.OK(t, err)
		test.Equals(t, expUpdate, actTask.LocalUpdate)
		test.Equals(t, task.STATUS_UPDATED, actTask.LocalStatus)
	})

	t.Run("markdispatched", func(t *testing.T) {
		mem := storage.NewMemory()
		test.OK(t, mem.SetTasks(tasks))
		lt, err := mem.FindById(task2.Id)
		test.OK(t, err)
		test.OK(t, mem.MarkDispatched(lt.LocalId))
		act, err := mem.FindById(task2.Id)
		test.OK(t, err)
		test.Equals(t, task.STATUS_DISPATCHED, act.LocalStatus)
	})

	t.Run("add", func(t *testing.T) {

		action := "action"
		project := "project"
		due := task.NewDate(2021, 9, 4)
		recur := task.Daily{Start: task.NewDate(2021, 9, 5)}
		mem := storage.NewMemory()
		expUpdate := &task.LocalUpdate{
			Fields:  []string{task.FIELD_ACTION, task.FIELD_PROJECT, task.FIELD_DUE, task.FIELD_RECUR},
			Action:  action,
			Project: project,
			Due:     due,
			Recur:   recur,
		}
		act1, err := mem.Add(expUpdate)
		test.OK(t, err)
		test.Assert(t, act1.Id != "", "id was empty")
		test.Equals(t, task.FOLDER_NEW, act1.Folder)
		test.Equals(t, "", act1.Action)
		test.Equals(t, "", act1.Project)
		test.Assert(t, act1.Due.IsZero(), "date was not zero")
		test.Equals(t, nil, act1.Recur)
		test.Equals(t, 0, act1.Version)
		test.Equals(t, 1, act1.LocalId)
		test.Equals(t, task.STATUS_UPDATED, act1.LocalStatus)
		test.Equals(t, expUpdate, act1.LocalUpdate)

		act2, err := mem.FindById(act1.Id)
		test.OK(t, err)
		test.Equals(t, act1, act2)
	})
}
