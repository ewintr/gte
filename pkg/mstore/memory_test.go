package mstore_test

import (
	"fmt"
	"sort"
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
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

func TestMemoryFolders(t *testing.T) {
	for _, tc := range []struct {
		name    string
		folders []string
	}{
		{
			name:    "one",
			folders: []string{"one"},
		},
		{
			name:    "many",
			folders: []string{"one", "two", "three"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem, err := mstore.NewMemory(tc.folders)
			test.OK(t, err)
			actFolders, err := mem.Folders()
			test.OK(t, err)
			expFolders := tc.folders
			sort.Strings(expFolders)
			test.Equals(t, expFolders, actFolders)
		})
	}
}

func TestMemorySelect(t *testing.T) {
	mem, err := mstore.NewMemory([]string{"one", "two", "three"})
	test.OK(t, err)
	for _, tc := range []struct {
		name   string
		folder string
		exp    error
	}{
		{
			name: "empty",
			exp:  mstore.ErrFolderDoesNotExist,
		},
		{
			name:   "not present",
			folder: "four",
			exp:    mstore.ErrFolderDoesNotExist,
		},
		{
			name:   "present",
			folder: "two",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, mem.Select(tc.folder))
			if tc.exp == nil {
				test.Equals(t, tc.folder, mem.Selected)
			}
		})
	}
}

func TestMemoryAdd(t *testing.T) {
	folder := "folder"

	for _, tc := range []struct {
		name    string
		subject string
		exp     error
	}{
		{
			name: "empty",
			exp:  mstore.ErrInvalidMessage,
		},
		{
			name:    "valid",
			subject: "subject",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem, err := mstore.NewMemory([]string{folder})
			test.OK(t, err)
			test.Equals(t, tc.exp, mem.Add(folder, tc.subject, ""))
		})
	}
}

func TestMemoryMessages(t *testing.T) {
	folderA := "folderA"
	folderB := "folderB"

	t.Run("no folder selected", func(t *testing.T) {
		mem, err := mstore.NewMemory([]string{folderA})
		test.OK(t, err)
		_, err = mem.Messages()
		test.Equals(t, mstore.ErrNoFolderSelected, err)
	})

	for _, tc := range []struct {
		name   string
		folder string
		amount int
		expErr error
	}{
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
					Subject: fmt.Sprintf("subject-%d", i),
					Body:    fmt.Sprintf("body-%d", i),
				}
				if tc.folder == folderA {
					expMessages = append(expMessages, m)
				}
				test.OK(t, mem.Add(folderA, m.Subject, m.Body))
			}

			test.OK(t, mem.Select(tc.folder))
			actMessages, err := mem.Messages()
			test.Equals(t, tc.expErr, err)
			test.Equals(t, expMessages, actMessages)
		})
	}
}

func TestMemoryRemove(t *testing.T) {
	folderA, folderB := "folderA", "folderB"

	t.Run("no folder selected", func(t *testing.T) {
		mem, err := mstore.NewMemory([]string{folderA})
		test.OK(t, err)
		test.Equals(t, mstore.ErrNoFolderSelected, mem.Remove(uint32(3)))
	})

	mem, err := mstore.NewMemory([]string{folderA, folderB})
	test.OK(t, err)
	for i := 1; i <= 3; i++ {
		test.OK(t, mem.Add(folderA, fmt.Sprintf("subject-%d", i), ""))
	}
	for _, tc := range []struct {
		name    string
		folder  string
		uid     uint32
		expUids []uint32
		expErr  error
	}{
		{
			name:   "invalid uid",
			folder: folderA,
			uid:    uint32(0),
			expErr: mstore.ErrInvalidUid,
		},
		{
			name:   "empty",
			folder: folderB,
			uid:    uint32(1),
			expErr: mstore.ErrMessageDoesNotExist,
		},
		{
			name:    "valid",
			folder:  folderA,
			uid:     uint32(2),
			expUids: []uint32{1, 3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.OK(t, mem.Select(tc.folder))
			test.Equals(t, tc.expErr, mem.Remove(tc.uid))
			if tc.expErr != nil {
				return
			}
			actUids := []uint32{}
			actMsgs, err := mem.Messages()
			test.OK(t, err)
			for _, m := range actMsgs {
				actUids = append(actUids, m.Uid)
			}
			test.Equals(t, tc.expUids, actUids)
		})
	}
}
