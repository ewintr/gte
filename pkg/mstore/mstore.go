package mstore

type Folder struct {
	Name    string
	Version uint32
}

type Message struct {
	ID      uint32
	Subject string
	Body    string
}

type MStorer interface {
	Folders() ([]string, error)
	Select(folder string) error
	Messages() ([]*Message, error)
	Add(message *Message) error
	Remove(id string) error
}
