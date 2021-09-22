package mstore

import (
	"errors"
)

var (
	ErrFolderDoesNotExist  = errors.New("folder does not exist")
	ErrNoFolderSelected    = errors.New("no folder selected")
	ErrInvalidUid          = errors.New("invalid uid")
	ErrMessageDoesNotExist = errors.New("message does not exist")
	ErrInvalidMessage      = errors.New("message is invalid")
)

type Message struct {
	Uid     uint32
	Folder  string
	Subject string
	Body    string
}

func (m *Message) Valid() bool {
	return m.Uid != 0 && m.Subject != "" && m.Folder != ""
}

func (m *Message) Equal(n *Message) bool {
	if m.Folder != n.Folder {
		return false
	}
	if m.Subject != n.Subject {
		return false
	}
	if m.Body != n.Body {
		return false
	}

	return true
}

type MStorer interface {
	Messages(folder string) ([]*Message, error)
	Add(folder, subject, body string) error
	Remove(msg *Message) error
}
