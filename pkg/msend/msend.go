package msend

import "errors"

var (
	ErrSendFail = errors.New("could not send message")
)

type Message struct {
	Subject string
	Body    string
}

type MSender interface {
	Send(msg *Message) error
}
