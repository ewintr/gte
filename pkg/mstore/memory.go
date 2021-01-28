package mstore

type Memory struct {
	messages map[*Folder][]*Message
}
