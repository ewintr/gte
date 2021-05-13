package msend

type Memory struct {
	Messages []*Message
}

func NewMemory() *Memory {
	return &Memory{
		Messages: []*Message{},
	}
}

func (mem *Memory) Send(msg *Message) error {
	mem.Messages = append(mem.Messages, msg)

	return nil
}
