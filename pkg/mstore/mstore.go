package mstore

import "errors"

var (
	ErrFolderDoesNotExist  = errors.New("folder does not exist")
	ErrNoFolderSelected    = errors.New("no folder selected")
	ErrInvalidUid          = errors.New("invalid uid")
	ErrMessageDoesNotExist = errors.New("message does not exist")
	ErrInvalidMessage      = errors.New("message is invalid")
)

type Message struct {
	Uid     uint32
	Subject string
	Body    string
}

func (m *Message) Valid() bool {
	return m.Uid != 0 && m.Subject != ""
}

type MStorer interface {
	Folders() ([]string, error)
	Select(folder string) error
	Messages() ([]*Message, error)
	Add(message *Message) error
	Remove(uid uint32) error
}
