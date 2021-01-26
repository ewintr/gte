package mstore

import (
	"github.com/emersion/go-imap/client"
)

type EMailStoreConfiguration struct {
	IMAPURL      string
	IMAPPort     int
	IMAPUsername string
	IMAPPassword string
	SMTPURL      string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type EMailStore struct {
	imapClient *client.Client
}

func NewEmailStore(config *EMailStoreConfiguration) (*EMailStore, error) {

	return &EMailStore{}, nil
}
