package msend

type Memory struct {
	Messages []*Message
	Err      error
}

func NewMemory() *Memory {
	return &Memory{
		Messages: []*Message{},
	}
}

func (mem *Memory) Send(msg *Message) error {
	if mem.Err != nil {
		return mem.Err
	}

	mem.Messages = append(mem.Messages, msg)

	return nil
}
