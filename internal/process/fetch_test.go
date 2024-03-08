package process_test

import (
	"sort"
	"testing"

	"code.ewintr.nl/go-kit/test"
	"code.ewintr.nl/gte/internal/process"
	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/internal/task"
	"code.ewintr.nl/gte/pkg/mstore"
)

func TestFetchProcess(t *testing.T) {
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
		Folder:  task.FOLDER_UNPLANNED,
	}
	task3 := &task.Task{
		Id:      "id3",
		Version: 1,
		Action:  "action3",
		Folder:  task.FOLDER_PLANNED,
	}

	localTask1 := &task.LocalTask{Task: *task1, LocalUpdate: &task.LocalUpdate{}, LocalStatus: task.STATUS_FETCHED}
	localTask2 := &task.LocalTask{Task: *task2, LocalUpdate: &task.LocalUpdate{}, LocalStatus: task.STATUS_FETCHED}
	localTask3 := &task.LocalTask{Task: *task3, LocalUpdate: &task.LocalUpdate{}, LocalStatus: task.STATUS_FETCHED}

	mstorer, err := mstore.NewMemory(task.KnownFolders)
	test.OK(t, err)
	test.OK(t, mstorer.Add(task1.Folder, task1.FormatSubject(), task1.FormatBody()))
	test.OK(t, mstorer.Add(task2.Folder, task2.FormatSubject(), task2.FormatBody()))
	test.OK(t, mstorer.Add(task3.Folder, task3.FormatSubject(), task3.FormatBody()))
	remote := storage.NewRemoteRepository(mstorer)
	local := storage.NewMemory()

	t.Run("all", func(t *testing.T) {
		syncer := process.NewFetch(remote, local)
		actResult, err := syncer.Process()
		test.OK(t, err)
		test.Equals(t, 3, actResult.Count)
		actTasks, err := local.FindAll()
		test.OK(t, err)
		for _, a := range actTasks {
			a.LocalId = 0
			a.Message = nil
		}
		exp := task.ById([]*task.LocalTask{localTask1, localTask2, localTask3})
		sExp := task.ById(exp)
		sAct := task.ById(actTasks)
		sort.Sort(sAct)
		sort.Sort(sExp)
		test.Equals(t, sExp, sAct)
	})

	t.Run("planned", func(t *testing.T) {
		syncer := process.NewFetch(remote, local, task.FOLDER_PLANNED)
		actResult, err := syncer.Process()
		test.OK(t, err)
		test.Equals(t, 1, actResult.Count)
		actTasks, err := local.FindAll()
		test.OK(t, err)
		for _, a := range actTasks {
			a.LocalId = 0
			a.Message = nil
		}
		exp := task.ById([]*task.LocalTask{localTask3})
		sExp := task.ById(exp)
		sAct := task.ById(actTasks)
		sort.Sort(sAct)
		sort.Sort(sExp)
		test.Equals(t, sExp, sAct)
	})

}
