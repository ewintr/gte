package process_test

import (
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
	"ewintr.nl/gte/pkg/mstore"
)

func TestInboxProcess(t *testing.T) {
	for _, tc := range []struct {
		name     string
		messages map[string][]*mstore.Message
		expCount int
		expMsgs  map[string][]*mstore.Message
	}{
		{
			name: "empty",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {},
			},
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {},
			},
		},
		{
			name: "all flavors",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {
					{
						Subject: "to new",
					},
					{
						Subject: "to recurring",
						Body:    "recur: 2021-05-14, daily\nid: xxx-xxx-a\nversion: 1",
					},
					{
						Subject: "to planned",
						Body:    "due: 2021-05-14\nid: xxx-xxx-b\nversion: 1",
					},
					{
						Subject: "to unplanned",
						Body:    "id: xxx-xxx-c\nversion: 1",
					},
				},
			},
			expCount: 4,
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX:     {},
				task.FOLDER_NEW:       {{Subject: "to new"}},
				task.FOLDER_RECURRING: {{Subject: "to recurring"}},
				task.FOLDER_PLANNED:   {{Subject: "2021-05-14 (friday) - to planned"}},
				task.FOLDER_UNPLANNED: {{Subject: "to unplanned"}},
			},
		},
		{
			name: "cleanup",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {{
					Subject: "new version",
					Body:    "id: xxx-xxx\nversion: 3",
				}},
				task.FOLDER_UNPLANNED: {{
					Subject: "old version",
					Body:    "id: xxx-xxx\nversion: 3",
				}},
			},
			expCount: 1,
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX:     {},
				task.FOLDER_UNPLANNED: {{Subject: "new version"}},
			},
		},
		{
			name: "cleanup version conflict",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {{
					Subject: "new version",
					Body:    "id: xxx-xxx\nversion: 3",
				}},
				task.FOLDER_UNPLANNED: {{
					Subject: "not really old version",
					Body:    "id: xxx-xxx\nversion: 5",
				}},
			},
			expCount: 1,
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX:     {},
				task.FOLDER_UNPLANNED: {{Subject: "not really old version"}},
			},
		},
		{
			name: "remove done",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {{
					Subject: "is done",
					Body:    "id: xxx-xxx\nversion: 1\ndone: true\n",
				}},
				task.FOLDER_UNPLANNED: {{
					Subject: "the task",
					Body:    "id: xxx-xxx\nversion: 1\n",
				}},
			},
			expCount: 1,
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX:     {},
				task.FOLDER_UNPLANNED: {},
			},
		},
		{
			name: "deduplicate",
			messages: map[string][]*mstore.Message{
				task.FOLDER_INBOX: {
					{
						Subject: "version 2",
						Body:    "id: xxx-xxx\nversion: 1\n",
					},
					{
						Subject: "version 2b",
						Body:    "id: xxx-xxx\nversion: 1\n",
					},
				},
				task.FOLDER_UNPLANNED: {{
					Subject: "the task",
					Body:    "id: xxx-xxx\nversion: 1\n",
				}},
			},
			expCount: 1,
			expMsgs: map[string][]*mstore.Message{
				task.FOLDER_INBOX:     {},
				task.FOLDER_UNPLANNED: {{Subject: "version 2b"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mstorer, err := mstore.NewMemory([]string{
				task.FOLDER_INBOX,
				task.FOLDER_NEW,
				task.FOLDER_RECURRING,
				task.FOLDER_PLANNED,
				task.FOLDER_UNPLANNED,
			})
			test.OK(t, err)
			for folder, messages := range tc.messages {
				for _, m := range messages {
					test.OK(t, mstorer.Add(folder, m.Subject, m.Body))
				}
			}

			inboxProc := process.NewInbox(storage.NewRemoteRepository(mstorer))
			actResult, err := inboxProc.Process()

			test.OK(t, err)
			test.Equals(t, tc.expCount, actResult.Count)
			for folder, expMessages := range tc.expMsgs {
				actMessages, err := mstorer.Messages(folder)
				test.OK(t, err)
				test.Equals(t, len(expMessages), len(actMessages))
				if len(expMessages) == 0 {

					continue
				}
				test.Equals(t, expMessages[0].Subject, actMessages[0].Subject)
			}
		})
	}
}
