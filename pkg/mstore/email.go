package mstore

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
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

type ImapConfiguration struct {
	ImapUrl      string
	ImapUsername string
	ImapPassword string
}

func (esc *ImapConfiguration) Valid() bool {
	if esc.ImapUrl == "" {
		return false
	}
	if esc.ImapUsername == "" || esc.ImapPassword == "" {
		return false
	}

	return true
}

type Imap struct {
	imap       *client.Client
	mboxStatus *imap.MailboxStatus
}

func ImapConnect(conf *ImapConfiguration) (*Imap, error) {
	imap, err := client.DialTLS(conf.ImapUrl, nil)
	if err != nil {
		return &Imap{}, err
	}
	if err := imap.Login(conf.ImapUsername, conf.ImapPassword); err != nil {
		return &Imap{}, err
	}

	return &Imap{
		imap: imap,
	}, nil
}

func (es *Imap) Disconnect() {
	es.imap.Logout()
}

func (es *Imap) Folders() ([]string, error) {
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

func (es *Imap) selectFolder(folder string) error {
	status, err := es.imap.Select(folder, false)
	if err != nil {
		return err
	}

	es.mboxStatus = status

	return nil
}

func (es *Imap) Messages(folder string) ([]*Message, error) {
	if err := es.selectFolder(folder); err != nil {
		return []*Message{}, err
	}

	if es.mboxStatus.Messages == 0 {
		return []*Message{}, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), es.mboxStatus.Messages)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchUid, section.FetchItem()}

	imsg, done := make(chan *imap.Message), make(chan error)
	go func() {
		done <- es.imap.Fetch(seqset, items, imsg)
	}()

	messages := []*Message{}
	for m := range imsg {
		r := m.GetBody(section)
		if r == nil {
			return []*Message{}, fmt.Errorf("server didn't returned message body")
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			return []*Message{}, err
		}

		// Print some info about the message
		header := mr.Header
		subject, err := header.Subject()
		if err != nil {
			return []*Message{}, err
		}

		// Process each message's part
		body := []byte(``)
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				return []*Message{}, err
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				body, _ = ioutil.ReadAll(p.Body)
			}
		}

		messages = append(messages, &Message{
			Uid:     m.Uid,
			Folder:  folder,
			Subject: subject,
			Body:    string(body),
		})
	}

	if err := <-done; err != nil {
		return []*Message{}, err
	}

	return messages, nil
}

func (es *Imap) Add(folder, subject, body string) error {
	msgStr := fmt.Sprintf(`From: todo <process@erikwinter.nl>
Subject: %s

%s`, subject, body)

	msg := NewBody(msgStr)

	return es.imap.Append(folder, nil, time.Time{}, imap.Literal(msg))
}

func (es *Imap) Remove(msg *Message) error {
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
