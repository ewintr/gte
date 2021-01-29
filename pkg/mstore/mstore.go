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
	Folder  string
	Subject string
	Body    string
}

func (m *Message) Valid() bool {
	return m.Uid != 0 && m.Subject != "" && m.Folder != ""
}

type MStorer interface {
	Folders() ([]string, error)
	Messages(folder string) ([]*Message, error)
	Add(folder, subject, body string) error
	Remove(msg *Message) error
}
