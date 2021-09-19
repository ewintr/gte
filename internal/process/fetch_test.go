package process_test

import (
	"sort"
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
	"ewintr.nl/gte/pkg/mstore"
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

	localTask1 := &task.LocalTask{Task: *task1, LocalUpdate: &task.LocalUpdate{}, LocalStatus: task.STATUS_FETCHED}
	localTask2 := &task.LocalTask{Task: *task2, LocalUpdate: &task.LocalUpdate{}, LocalStatus: task.STATUS_FETCHED}

	mstorer, err := mstore.NewMemory(task.KnownFolders)
	test.OK(t, err)
	test.OK(t, mstorer.Add(task1.Folder, task1.FormatSubject(), task1.FormatBody()))
	test.OK(t, mstorer.Add(task2.Folder, task2.FormatSubject(), task2.FormatBody()))
	remote := storage.NewRemoteRepository(mstorer)
	local := storage.NewMemory()

	syncer := process.NewFetch(remote, local)
	actResult, err := syncer.Process()
	test.OK(t, err)
	test.Equals(t, 2, actResult.Count)
	actTasks, err := local.FindAll()
	test.OK(t, err)
	for _, a := range actTasks {
		a.LocalId = 0
		a.Message = nil
	}
	exp := task.ById([]*task.LocalTask{localTask1, localTask2})
	sExp := task.ById(exp)
	sAct := task.ById(actTasks)
	sort.Sort(sAct)
	sort.Sort(sExp)
	test.Equals(t, sExp, sAct)
}
