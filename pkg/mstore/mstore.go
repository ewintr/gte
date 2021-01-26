package mstore

type Folder struct {
	Name string
}

type Message struct {
	Subject string
}

type MStorer interface {
	Folders() ([]Folder, error)
	Messages(folder Folder) ([]Message, error)
	Move(message Message, folder Folder) error
	Update(message Message) error
	Add(message Message) error
}
