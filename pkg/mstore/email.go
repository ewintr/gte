package mstore

import (
	"fmt"

	"github.com/emersion/go-imap"
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

func (esc *EMailStoreConfiguration) Valid() bool {
	if esc.IMAPURL == "" || esc.IMAPPort == 0 {
		return false
	}
	if esc.IMAPUsername == "" || esc.IMAPPassword == "" {
		return false
	}

	return true
}

type EMailStore struct {
	imapClient *client.Client
}

func EMailConnect(conf *EMailStoreConfiguration) (*EMailStore, error) {
	imapClient, err := client.DialTLS(fmt.Sprintf("%s:%d", conf.IMAPURL, conf.IMAPPort), nil)
	if err != nil {
		return &EMailStore{}, err
	}
	if err := imapClient.Login(conf.IMAPUsername, conf.IMAPPassword); err != nil {
		return &EMailStore{}, err
	}

	return &EMailStore{
		imapClient: imapClient,
	}, nil
}

func (es EMailStore) Disconnect() {
	es.imapClient.Logout()
}

func (es EMailStore) Folders() ([]*Folder, error) {
	boxes, done := make(chan *imap.MailboxInfo), make(chan error)
	go func() {
		done <- es.imapClient.List("", "*", boxes)
	}()

	folders := []*Folder{}
	for b := range boxes {
		folders = append(folders, &Folder{
			Name: b.Name,
		})
	}

	if err := <-done; err != nil {
		return []*Folder{}, err
	}

	return folders, nil
}
