package task_test

import (
	"errors"
	"fmt"
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/internal/task"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
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
			expErr: task.ErrMStoreError,
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
			repo := task.NewRepository(store)
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
	oldFolder := "old folder"
	folder := "folder"
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
			expErr: task.ErrInvalidTask,
		},
		{
			name: "task without message",
			task: &task.Task{
				Id:      id,
				Folder:  folder,
				Action:  action,
				Current: true,
				Dirty:   true,
			},
			expErr: task.ErrMStoreError,
		},
		{
			name: "outdated task",
			task: &task.Task{
				Id:     id,
				Folder: folder,
				Action: action,
				Dirty:  true,
			},
			expErr: task.ErrOutdatedTask,
		},
		/*
			{
				name: "unchanged task",
				task: &task.Task{
					Id:      id,
					Folder:  folder,
					Action:  action,
					Current: true,
				},
				expMsgs: []*mstore.Message{},
			},
		*/
		{
			name: "changed task",
			task: &task.Task{
				Id:      id,
				Folder:  folder,
				Action:  action,
				Current: true,
				Dirty:   true,
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

			repo := task.NewRepository(mem)

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
	folderNew := "New"
	folderPlanned := "Planned"
	folders := []string{"INBOX", folderNew, "Recurring",
		folderPlanned, "Unplanned",
	}
	id := "id"
	subject := "subject"

	mem, err := mstore.NewMemory(folders)
	test.OK(t, err)

	for v := 1; v <= 3; v++ {
		body := fmt.Sprintf(`
id: %s
version: %d
`, id, v)
		folder := folderNew
		if v%2 == 1 {
			folder = folderPlanned
		}
		test.OK(t, mem.Add(folder, subject, body))
	}

	repo := task.NewRepository(mem)
	test.OK(t, repo.CleanUp())

	expNew := []*mstore.Message{}
	actNew, err := mem.Messages(folderNew)
	test.OK(t, err)
	test.Equals(t, expNew, actNew)
	expPlanned := []*mstore.Message{{
		Uid:     3,
		Folder:  folderPlanned,
		Subject: subject,
		Body: `
id: id
version: 3
`,
	}}
	actPlanned, err := mem.Messages(folderPlanned)
	test.OK(t, err)
	test.Equals(t, expPlanned, actPlanned)
}