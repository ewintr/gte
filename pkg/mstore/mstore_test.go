package mstore_test

import (
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

func TestMessageValid(t *testing.T) {
	for _, tc := range []struct {
		name    string
		message *mstore.Message
		exp     bool
	}{
		{
			name:    "empty",
			message: &mstore.Message{},
		},
		{
			name: "no uid",
			message: &mstore.Message{
				Subject: "subject",
				Folder:  "folder",
				Body:    "body",
			},
		},
		{
			name: "no  folder",
			message: &mstore.Message{
				Uid:     1,
				Subject: "subject",
				Body:    "body",
			},
		},
		{
			name: "no subject",
			message: &mstore.Message{
				Uid:    1,
				Folder: "folder",
				Body:   "body",
			},
		},
		{
			name: "no body",
			message: &mstore.Message{
				Uid:     1,
				Folder:  "folder",
				Subject: "subject",
			},
			exp: true,
		},
		{
			name: "all present",
			message: &mstore.Message{
				Uid:     1,
				Folder:  "folder",
				Subject: "subject",
				Body:    "body",
			},
			exp: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, tc.message.Valid())
		})
	}
}
