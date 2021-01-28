package mstore

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

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

func (es *Email) FolderNames() ([]string, error) {
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

func (es *Email) Select(folder string) error {
	status, err := es.imap.Select(folder, false)
	if err != nil {
		return err
	}
	fmt.Printf("status: %+v\n", status)

	es.mboxStatus = status

	return nil
}

func (es *Email) Messages() ([]*Message, error) {
	if es.mboxStatus == nil {
		return []*Message{}, fmt.Errorf("no mailbox selected")
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
		//fmt.Printf("%+v\n", m)
		messages = append(messages, &Message{
			Uid:     m.Uid,
			Subject: m.Envelope.Subject,
		})
	}

	if err := <-done; err != nil {
		return []*Message{}, err
	}

	return messages, nil
}

func (es *Email) Append(mbox string, msg imap.Literal) error {
	return es.imap.Append(mbox, nil, time.Time{}, msg)
}

func (es *Email) Remove(uid uint32) error {
	if uid == 0 {
		return fmt.Errorf("invalid uid: %d", uid)
	}
	if es.mboxStatus == nil {
		return fmt.Errorf("no mailbox selected")
	}

	// set deleted flag
	seqset := new(imap.SeqSet)
	seqset.AddRange(uid, uid)
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
