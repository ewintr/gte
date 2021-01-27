package mstore

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type EMailStoreConfiguration struct {
	IMAPURL      string
	IMAPUsername string
	IMAPPassword string
}

func (esc *EMailStoreConfiguration) Valid() bool {
	if esc.IMAPURL == "" {
		return false
	}
	if esc.IMAPUsername == "" || esc.IMAPPassword == "" {
		return false
	}

	return true
}

type smtpConf struct {
	url string
	to  string
}

type EMailStore struct {
	imap *client.Client
}

func EMailConnect(conf *EMailStoreConfiguration) (*EMailStore, error) {
	imap, err := client.DialTLS(conf.IMAPURL, nil)
	if err != nil {
		return &EMailStore{}, err
	}
	if err := imap.Login(conf.IMAPUsername, conf.IMAPPassword); err != nil {
		return &EMailStore{}, err
	}

	return &EMailStore{
		imap: imap,
	}, nil
}

func (es *EMailStore) Disconnect() {
	es.imap.Logout()
}

func (es *EMailStore) Folders() ([]*Folder, error) {
	boxes, done := make(chan *imap.MailboxInfo), make(chan error)
	go func() {
		done <- es.imap.List("", "*", boxes)
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

func (es *EMailStore) Inbox() ([]*Message, error) {
	mbox, err := es.imap.Select("INBOX", false)
	if err != nil {
		return []*Message{}, err
	}
	fmt.Println("Flags for INBOX:", mbox.Flags)

	fmt.Println("Messages: ", mbox.Messages)

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), mbox.Messages)

	imsg, done := make(chan *imap.Message), make(chan error)
	go func() {
		done <- es.imap.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, imsg)
	}()

	messages := []*Message{}
	for m := range imsg {
		messages = append(messages, &Message{
			Subject: m.Envelope.Subject,
		})
	}

	if err := <-done; err != nil {
		return []*Message{}, err
	}

	return messages, nil
}

func (es *EMailStore) Append(mbox string, msg imap.Literal) error {
	return es.imap.Append(mbox, nil, time.Time{}, msg)
}
