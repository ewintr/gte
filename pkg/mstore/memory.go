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
	nextUid  uint32
	folders  []string
	messages map[string][]*Message
}

func NewMemory(folders []string) (*Memory, error) {
	if len(folders) == 0 {
		return &Memory{}, ErrInvalidFolderSet
	}
	sort.Strings(folders)

	msg := make(map[string][]*Message)
	for _, f := range folders {
		if f == "" {
			return &Memory{}, ErrInvalidFolderName
		}
		msg[f] = []*Message{}
	}

	return &Memory{
		messages: msg,
		folders:  folders,
		nextUid:  uint32(1),
	}, nil
}

func (mem *Memory) Folders() ([]string, error) {
	return mem.folders, nil
}

func (mem *Memory) Add(folder, subject, body string) error {
	if subject == "" {
		return ErrInvalidMessage
	}
	if _, ok := mem.messages[folder]; !ok {
		return ErrFolderDoesNotExist
	}

	mem.messages[folder] = append(mem.messages[folder], &Message{
		Uid:     mem.nextUid,
		Folder:  folder,
		Subject: subject,
		Body:    body,
	})
	mem.nextUid++

	return nil
}

func (mem *Memory) Messages(folder string) ([]*Message, error) {
	if _, ok := mem.messages[folder]; !ok {
		return []*Message{}, ErrFolderDoesNotExist
	}

	return mem.messages[folder], nil
}

func (mem *Memory) Remove(msg *Message) error {
	if msg == nil || !msg.Valid() {
		return ErrInvalidMessage
	}
	if _, ok := mem.messages[msg.Folder]; !ok {
		return ErrFolderDoesNotExist
	}

	for i, m := range mem.messages[msg.Folder] {
		if m.Uid == msg.Uid {
			mem.messages[msg.Folder] = append(mem.messages[msg.Folder][:i], mem.messages[msg.Folder][i+1:]...)

			return nil
		}
	}

	return ErrMessageDoesNotExist
}
