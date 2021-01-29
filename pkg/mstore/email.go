package mstore

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Body struct {
	reader io.Reader
	length int
}

func NewBody(msg string) *Body {

	return &Body{
		reader: strings.NewReader(msg),
		length: len([]byte(msg)),
	}
}

func (b *Body) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

func (b *Body) Len() int {
	return b.length
}

type EmailConfiguration struct {
	IMAPURL      string
	IMAPUsername string
	IMAPPassword string
}

func (esc *EmailConfiguration) Valid() bool {
	if esc.IMAPURL == "" {
		return false
	}
	if esc.IMAPUsername == "" || esc.IMAPPassword == "" {
		return false
	}

	return true
}

type Email struct {
	imap       *client.Client
	mboxStatus *imap.MailboxStatus
}

func EmailConnect(conf *EmailConfiguration) (*Email, error) {
	imap, err := client.DialTLS(conf.IMAPURL, nil)
	if err != nil {
		return &Email{}, err
	}
	if err := imap.Login(conf.IMAPUsername, conf.IMAPPassword); err != nil {
		return &Email{}, err
	}

	return &Email{
		imap: imap,
	}, nil
}

func (es *Email) Disconnect() {
	es.imap.Logout()
}

func (es *Email) Folders() ([]string, error) {
	boxes, done := make(chan *imap.MailboxInfo), make(chan error)
	go func() {
		done <- es.imap.List("", "*", boxes)
	}()

	folders := []string{}
	for b := range boxes {
		folders = append(folders, b.Name)
	}

	if err := <-done; err != nil {
		return []string{}, err
	}

	return folders, nil
}

func (es *Email) selectFolder(folder string) error {
	status, err := es.imap.Select(folder, false)
	if err != nil {
		return err
	}

	es.mboxStatus = status

	return nil
}

func (es *Email) Messages(folder string) ([]*Message, error) {
	if err := es.selectFolder(folder); err != nil {
		return []*Message{}, err
	}

	if es.mboxStatus.Messages == 0 {
		return []*Message{}, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), es.mboxStatus.Messages)

	imsg, done := make(chan *imap.Message), make(chan error)
	go func() {
		done <- es.imap.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, imsg)
	}()

	messages := []*Message{}
	for m := range imsg {
		messages = append(messages, &Message{
			Uid:     m.Uid,
			Folder:  folder,
			Subject: m.Envelope.Subject,
		})
	}

	if err := <-done; err != nil {
		return []*Message{}, err
	}

	return messages, nil
}

func (es *Email) Add(folder, subject, body string) error {
	msgStr := fmt.Sprintf(`From: todo <process@erikwinter.nl>
Subject: %s

%s`, subject, body)

	msg := NewBody(msgStr)

	return es.imap.Append(folder, nil, time.Time{}, imap.Literal(msg))
}

func (es *Email) Remove(msg *Message) error {
	if !msg.Valid() {
		return ErrInvalidMessage
	}

	if err := es.selectFolder(msg.Folder); err != nil {
		return err
	}

	// set deleted flag
	seqset := new(imap.SeqSet)
	seqset.AddRange(msg.Uid, msg.Uid)
	storeItem := imap.FormatFlagsOp(imap.SetFlags, true)
	err := es.imap.UidStore(seqset, storeItem, imap.FormatStringList([]string{imap.DeletedFlag}), nil)
	if err != nil {
		return err
	}

	// expunge box
	if err := es.imap.Expunge(nil); err != nil {
		return err
	}

	return nil
}
