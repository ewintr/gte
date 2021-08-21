package mstore

import (
	"errors"
	"fmt"
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
	var prt bool
	if m.Uid == 156 && n.Uid == 155 {
		prt = true
	}
	if m.Uid == 155 && n.Uid == 156 {
		prt = true
	}
	if m.Folder != n.Folder {
		if prt {
			fmt.Println("folder")
		}
		return false
	}
	if m.Subject != n.Subject {
		if prt {
			fmt.Println("subject")
		}
		return false
	}
	if m.Body != n.Body {
		if prt {
			fmt.Println("body")
		}
		return false
	}

	return true
}

type MStorer interface {
	Folders() ([]string, error)
	Messages(folder string) ([]*Message, error)
	Add(folder, subject, body string) error
	Remove(msg *Message) error
}
