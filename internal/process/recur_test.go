package process_test

import (
	"fmt"
	"testing"
	"time"

	"go-mod.ewintr.nl/go-kit/test"
	"go-mod.ewintr.nl/gte/internal/process"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/internal/task"
	"go-mod.ewintr.nl/gte/pkg/msend"
	"go-mod.ewintr.nl/gte/pkg/mstore"
)

func TestRecurProcess(t *testing.T) {
	strFormat := "2006-01-02"
	todayStr := time.Now().Format(strFormat)
	nextMonthStr := time.Now().Add(30 * 24 * time.Hour).Format(strFormat)
	tomorrowStr := task.Today().Add(1).String()
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
					Subject: "recurring",
					Body:    fmt.Sprintf("recur: %s, daily\nid: xxx-xxx\nversion: 1", todayStr),
				},
				{
					Subject: "not recurring",
					Body:    fmt.Sprintf("recur: %s, daily\nid: xxx-xxx\nversion: 1", nextMonthStr),
				},
			},
			expCount: 1,
			expMsgs: []*msend.Message{
				{Subject: fmt.Sprintf("%s - recurring", tomorrowStr)},
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

			recurProc := process.NewRecur(storage.NewRemoteRepository(mstorer), storage.NewDispatcher(msender), 1)
			actResult, err := recurProc.Process()
			test.OK(t, err)
			test.Equals(t, tc.expCount, actResult.Count)
			for i, expMsg := range tc.expMsgs {
				test.Equals(t, expMsg.Subject, msender.Messages[i].Subject)
			}
		})
	}
}
