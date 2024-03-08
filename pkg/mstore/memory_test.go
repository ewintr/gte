package mstore_test

import (
	"fmt"
	"testing"

	"code.ewintr.nl/go-kit/test"
	"code.ewintr.nl/gte/pkg/mstore"
)

func TestNewMemory(t *testing.T) {
	for _, tc := range []struct {
		name    string
		folders []string
		exp     error
	}{
		{
			name:    "empty",
			folders: []string{},
			exp:     mstore.ErrInvalidFolderSet,
		},
		{
			name: "nil",
			exp:  mstore.ErrInvalidFolderSet,
		},
		{
			name:    "empty string",
			folders: []string{""},
			exp:     mstore.ErrInvalidFolderName,
		},
		{
			name:    "one",
			folders: []string{"one"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mstore.NewMemory(tc.folders)
			test.Equals(t, tc.exp, err)
		})
	}
}

func TestMemoryAdd(t *testing.T) {
	folder := "folder"
	subject := "subject"

	for _, tc := range []struct {
		name    string
		folder  string
		subject string
		expMsgs []*mstore.Message
		expErr  error
	}{
		{
			name:   "empty",
			folder: folder,
			expErr: mstore.ErrInvalidMessage,
		},
		{
			name:    "invalid folder",
			folder:  "not there",
			subject: subject,
			expErr:  mstore.ErrFolderDoesNotExist,
		},
		{
			name:    "valid",
			folder:  folder,
			subject: subject,
			expMsgs: []*mstore.Message{
				{Uid: 1, Folder: folder, Subject: subject},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem, err := mstore.NewMemory([]string{folder})
			test.OK(t, err)
			test.Equals(t, tc.expErr, mem.Add(tc.folder, tc.subject, ""))
		})
	}
}

func TestMemoryMessages(t *testing.T) {
	folderA := "folderA"
	folderB := "folderB"

	for _, tc := range []struct {
		name   string
		folder string
		amount int
		expErr error
	}{
		{
			name:   "unknown folder",
			folder: "not there",
			expErr: mstore.ErrFolderDoesNotExist,
		},
		{
			name:   "empty folder",
			folder: folderB,
			amount: 3,
		},
		{
			name:   "one",
			folder: folderA,
			amount: 1,
		},
		{
			name:   "many",
			folder: folderA,
			amount: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem, err := mstore.NewMemory([]string{folderA, folderB})
			test.OK(t, err)

			expMessages := []*mstore.Message{}
			for i := 1; i <= tc.amount; i++ {
				m := &mstore.Message{
					Uid:     uint32(i),
					Folder:  folderA,
					Subject: fmt.Sprintf("subject-%d", i),
					Body:    fmt.Sprintf("body-%d", i),
				}
				if tc.folder == folderA {
					expMessages = append(expMessages, m)
				}
				test.OK(t, mem.Add(folderA, m.Subject, m.Body))
			}

			actMessages, err := mem.Messages(tc.folder)
			test.Equals(t, tc.expErr, err)
			test.Equals(t, expMessages, actMessages)
		})
	}
}

func TestMemoryRemove(t *testing.T) {
	folderA, folderB := "folderA", "folderB"
	subject := "subject"

	mem, err := mstore.NewMemory([]string{folderA, folderB})
	test.OK(t, err)
	for i := 1; i <= 3; i++ {
		test.OK(t, mem.Add(folderA, fmt.Sprintf("subject-%d", i), ""))
	}
	for _, tc := range []struct {
		name    string
		msg     *mstore.Message
		expUids []uint32
		expErr  error
	}{
		{
			name: "empty",
			msg: &mstore.Message{
				Uid:     1,
				Folder:  folderB,
				Subject: subject,
			},
			expErr: mstore.ErrMessageDoesNotExist,
		},
		{
			name:   "nil message",
			expErr: mstore.ErrInvalidMessage,
		},
		{
			name: "unknown folder",
			msg: &mstore.Message{
				Uid:     1,
				Folder:  "unknown",
				Subject: subject,
			},
			expErr: mstore.ErrFolderDoesNotExist,
		},
		{
			name: "valid",
			msg: &mstore.Message{
				Uid:     2,
				Folder:  folderA,
				Subject: subject,
			},
			expUids: []uint32{1, 3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.expErr, mem.Remove(tc.msg))
			if tc.expErr != nil {
				return
			}
			actUids := []uint32{}
			actMsgs, err := mem.Messages(tc.msg.Folder)
			test.OK(t, err)
			for _, m := range actMsgs {
				actUids = append(actUids, m.Uid)
			}
			test.Equals(t, tc.expUids, actUids)
		})
	}
}
