package storage_test

import (
	"errors"
	"fmt"
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func TestRepoFindAll(t *testing.T) {
	folderA := "folderA"
	folderB := "folderB"

	type msgs struct {
		Folder  string
		Subject string
	}

	for _, tc := range []struct {
		name     string
		tasks    []msgs
		folder   string
		expTasks int
		expErr   error
	}{
		{
			name:   "empty",
			folder: folderA,
		},
		{
			name:   "unknown folder",
			folder: "unknown",
			expErr: storage.ErrMStoreError,
		},
		{
			name:   "not empty",
			folder: folderA,
			tasks: []msgs{
				{Folder: folderA, Subject: "sub-1"},
				{Folder: folderA, Subject: "sub-2"},
				{Folder: folderB, Subject: "sub-3"},
				{Folder: folderA, Subject: "sub-4"},
			},
			expTasks: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, err := mstore.NewMemory([]string{folderA, folderB})
			test.OK(t, err)
			for _, task := range tc.tasks {
				test.OK(t, store.Add(task.Folder, task.Subject, "body"))
			}
			repo := storage.NewRemoteRepository(store)
			actTasks, err := repo.FindAll(tc.folder)
			test.Equals(t, true, errors.Is(err, tc.expErr))
			if err != nil {
				return
			}
			test.Equals(t, tc.expTasks, len(actTasks))

		})
	}
}

func TestRepoUpdate(t *testing.T) {
	id := "id"
	oldFolder := task.FOLDER_INBOX
	folder := task.FOLDER_NEW
	action := "action"

	oldMsg := &mstore.Message{
		Uid:     1,
		Folder:  oldFolder,
		Subject: "old subject",
		Body:    "old body",
	}

	for _, tc := range []struct {
		name    string
		task    *task.Task
		expErr  error
		expMsgs []*mstore.Message
	}{
		{
			name:   "nil task",
			expErr: storage.ErrInvalidTask,
		},
		{
			name: "task without message",
			task: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
			},
			expErr: storage.ErrMStoreError,
		},
		{
			name: "changed task",
			task: &task.Task{
				Id:      id,
				Folder:  folder,
				Action:  action,
				Message: oldMsg,
			},
			expMsgs: []*mstore.Message{
				{Uid: 2, Folder: folder, Subject: action},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem, err := mstore.NewMemory([]string{folder, oldFolder})
			test.OK(t, err)
			test.OK(t, mem.Add(oldMsg.Folder, oldMsg.Subject, oldMsg.Body))

			repo := storage.NewRemoteRepository(mem)

			actErr := repo.Update(tc.task)
			test.Equals(t, true, errors.Is(actErr, tc.expErr))
			if tc.expErr != nil {
				return
			}

			actMsgs, err := mem.Messages(folder)
			test.OK(t, err)
			for i, _ := range actMsgs {
				actMsgs[i].Body = ""
			}
			test.Equals(t, tc.expMsgs, actMsgs)
		})
	}
}

func TestRepoCleanUp(t *testing.T) {
	id := "id"
	subject := "subject"

	mem, err := mstore.NewMemory(task.KnownFolders)
	test.OK(t, err)

	for v := 1; v <= 3; v++ {
		body := fmt.Sprintf(`
id: %s
version: %d
`, id, v)
		folder := task.FOLDER_NEW
		if v%2 == 1 {
			folder = task.FOLDER_PLANNED
		}
		test.OK(t, mem.Add(folder, subject, body))
	}

	repo := storage.NewRemoteRepository(mem)
	test.OK(t, repo.CleanUp())

	expNew := []*mstore.Message{}
	actNew, err := mem.Messages(task.FOLDER_NEW)
	test.OK(t, err)
	test.Equals(t, expNew, actNew)
	expPlanned := []*mstore.Message{{
		Uid:     3,
		Folder:  task.FOLDER_PLANNED,
		Subject: subject,
		Body: `
id: id
version: 3
`,
	}}
	actPlanned, err := mem.Messages(task.FOLDER_PLANNED)
	test.OK(t, err)
	test.Equals(t, expPlanned, actPlanned)
}

func TestRepoRemove(t *testing.T) {
	mem, err := mstore.NewMemory(task.KnownFolders)
	test.OK(t, err)

	for id := 1; id <= 3; id++ {
		test.OK(t, mem.Add(task.FOLDER_PLANNED, "action", fmt.Sprintf("id: id-%d\n", id)))
	}
	remote := storage.NewRemoteRepository(mem)
	tasks := []*task.Task{
		{Id: "id-1"},
		{Id: "id-3"},
	}
	test.OK(t, remote.Remove(tasks))
	actMsgs, err := mem.Messages(task.FOLDER_PLANNED)
	expMsgs := []*mstore.Message{{
		Uid:     2,
		Folder:  task.FOLDER_PLANNED,
		Subject: "action",
		Body:    "id: id-2\n",
	}}
	test.Equals(t, expMsgs, actMsgs)

}
