package process_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func TestSyncProcess(t *testing.T) {
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

	mstorer, err := mstore.NewMemory(task.KnownFolders)
	test.OK(t, err)
	test.OK(t, mstorer.Add(task1.Folder, task1.FormatSubject(), task1.FormatBody()))
	test.OK(t, mstorer.Add(task2.Folder, task2.FormatSubject(), task2.FormatBody()))
	remote := storage.NewRemoteRepository(mstorer)
	local := storage.NewMemory()

	syncer := process.NewSync(remote, local)
	actResult, err := syncer.Process()
	test.OK(t, err)
	test.Equals(t, 2, actResult.Count)
	actTasks1, err := local.FindAllInFolder(task.FOLDER_NEW)
	test.OK(t, err)
	test.Equals(t, []*task.Task{task1}, actTasks1)
	actTasks2, err := local.FindAllInFolder(task.FOLDER_UNPLANNED)
	test.OK(t, err)
	test.Equals(t, []*task.Task{task2}, actTasks2)
}
