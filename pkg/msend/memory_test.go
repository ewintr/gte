package msend_test

import (
	"fmt"
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/pkg/msend"
)

func TestMemorySend(t *testing.T) {
	mem := msend.NewMemory()
	test.Equals(t, []*msend.Message{}, mem.Messages)

	msg1 := &msend.Message{Subject: "sub1", Body: "body1"}
	test.OK(t, mem.Send(msg1))
	test.Equals(t, []*msend.Message{msg1}, mem.Messages)

	msg2 := &msend.Message{Subject: "sub2", Body: "body2"}
	test.OK(t, mem.Send(msg2))
	test.Equals(t, []*msend.Message{msg1, msg2}, mem.Messages)

	expErr := fmt.Errorf("oh no")
	mem.Err = expErr
	actErr := mem.Send(msg1)
	test.Equals(t, expErr, actErr)
}
