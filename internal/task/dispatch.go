package task

import "git.ewintr.nl/gte/pkg/msend"

type Dispatcher struct {
	msender msend.MSender
}

func NewDispatcher(msender msend.MSender) *Dispatcher {
	return &Dispatcher{
		msender: msender,
	}
}

func (d *Dispatcher) Dispatch(t *Task) error {
	return d.msender.Send(&msend.Message{
		Subject: t.FormatSubject(),
		Body:    t.FormatBody(),
	})
}
