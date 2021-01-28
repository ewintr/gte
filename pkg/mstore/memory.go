package mstore

import (
	"errors"
	"sort"
)

var (
	ErrInvalidFolderSet  = errors.New("invalid folder set")
	ErrInvalidFolderName = errors.New("invalid folder name")
)

type Memory struct {
	Selected string
	nextUid  uint32
	messages map[string][]*Message
}

func NewMemory(folders []string) (*Memory, error) {
	if len(folders) == 0 {
		return &Memory{}, ErrInvalidFolderSet
	}

	msg := make(map[string][]*Message)
	for _, f := range folders {
		if f == "" {
			return &Memory{}, ErrInvalidFolderName
		}
		msg[f] = []*Message{}
	}

	return &Memory{
		messages: msg,
		nextUid:  uint32(1),
	}, nil
}

func (mem *Memory) Folders() ([]string, error) {
	folders := []string{}
	for f, _ := range mem.messages {
		folders = append(folders, f)
	}

	sort.Strings(folders)

	return folders, nil
}

func (mem *Memory) Select(folder string) error {
	if _, ok := mem.messages[folder]; !ok {
		return ErrFolderDoesNotExist
	}

	mem.Selected = folder

	return nil
}

func (mem *Memory) Add(subject, body string) error {
	if subject == "" {
		return ErrInvalidMessage
	}
	if mem.Selected == "" {
		return ErrNoFolderSelected
	}

	mem.messages[mem.Selected] = append(mem.messages[mem.Selected], &Message{
		Uid:     mem.nextUid,
		Subject: subject,
		Body:    body,
	})
	mem.nextUid++

	return nil
}

func (mem *Memory) Messages() ([]*Message, error) {
	if mem.Selected == "" {
		return []*Message{}, ErrNoFolderSelected
	}

	return mem.messages[mem.Selected], nil
}

func (mem *Memory) Remove(uid uint32) error {
	if uid == uint32(0) {
		return ErrInvalidUid
	}
	if mem.Selected == "" {
		return ErrNoFolderSelected
	}

	for i, m := range mem.messages[mem.Selected] {
		if m.Uid == uid {
			mem.messages[mem.Selected] = append(mem.messages[mem.Selected][:i], mem.messages[mem.Selected][i+1:]...)

			return nil
		}
	}
	return ErrMessageDoesNotExist
}
