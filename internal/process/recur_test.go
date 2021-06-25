package process_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

func TestRecurProcess(t *testing.T) {
	task.Today = task.NewDate(2021, 5, 14)
	for _, tc := range []struct {
		name      string
		recurMsgs []*mstore.Message
		expCount  int
		expMsgs   []*msend.Message
	}{
		{
			name:    "empty",
			expMsgs: []*msend.Message{},
		},
		{
			name: "one of two recurring",
			recurMsgs: []*mstore.Message{
				{
					Subject: "not recurring",
					Body:    "recur: 2021-05-20, daily\nid: xxx-xxx\nversion: 1",
				},
				{
					Subject: "recurring",
					Body:    "recur: 2021-05-10, daily\nid: xxx-xxx\nversion: 1",
				},
			},
			expCount: 1,
			expMsgs: []*msend.Message{
				{Subject: "2021-05-15 (saturday) - recurring"},
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
			for _, m := range tc.recurMsgs {
				test.OK(t, mstorer.Add(task.FOLDER_RECURRING, m.Subject, m.Body))
			}
			msender := msend.NewMemory()

			recurProc := process.NewRecur(task.NewRemoteRepository(mstorer), task.NewDispatcher(msender), 1)
			actResult, err := recurProc.Process()
			test.OK(t, err)
			test.Equals(t, tc.expCount, actResult.Count)
			for i, expMsg := range tc.expMsgs {
				test.Equals(t, expMsg.Subject, msender.Messages[i].Subject)
			}
		})
	}
}
