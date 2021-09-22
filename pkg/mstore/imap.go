package mstore

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

var (
	ErrIMAPInvalidConfig = errors.New("invalid imap configuration")
	ErrIMAPConnFailure   = errors.New("could not connect with imap")
	ErrIMAPNotConnected  = errors.New("unable to perform, not connected to imap")
	ErrIMAPServerProblem = errors.New("imap server was unable to perform operation")
)

type IMAPBody struct {
	reader io.Reader
	length int
}

func NewIMAPBody(msg string) *IMAPBody {
	return &IMAPBody{
		reader: strings.NewReader(msg),
		length: len([]byte(msg)),
	}
}

func (b *IMAPBody) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

func (b *IMAPBody) Len() int {
	return b.length
}

type IMAPConfig struct {
	IMAPURL          string
	IMAPUsername     string
	IMAPPassword     string
	IMAPFolderPrefix string
}

func (esc *IMAPConfig) Valid() bool {
	if esc.IMAPURL == "" {
		return false
	}
	if esc.IMAPUsername == "" || esc.IMAPPassword == "" {
		return false
	}

	return true
}

type IMAP struct {
	config     *IMAPConfig
	connected  bool
	client     *client.Client
	mboxStatus *imap.MailboxStatus
}

func NewIMAP(config *IMAPConfig) *IMAP {
	return &IMAP{
		config: config,
	}
}

func (im *IMAP) Connect() error {
	if !im.config.Valid() {
		return ErrIMAPInvalidConfig
	}
	if im.connected {
		return nil
	}

	cl, err := client.DialTLS(im.config.IMAPURL, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIMAPConnFailure, err)
	}
	if err := cl.Login(im.config.IMAPUsername, im.config.IMAPPassword); err != nil {
		return fmt.Errorf("%w: %v", ErrIMAPConnFailure, err)
	}

	im.client = cl
	im.connected = true

	return nil
}

func (im *IMAP) Close() {
	im.client.Logout()
	im.connected = false
}

func (im *IMAP) selectFolder(folder string) error {
	if !im.connected {
		return ErrIMAPNotConnected
	}

	folder = fmt.Sprintf("%s%s", im.config.IMAPFolderPrefix, folder)
	status, err := im.client.Select(folder, false)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrIMAPServerProblem, err)
	}

	im.mboxStatus = status

	return nil
}

func (im *IMAP) Messages(folder string) ([]*Message, error) {
	if err := im.Connect(); err != nil {
		return []*Message{}, err
	}
	defer im.Close()

	if err := im.selectFolder(folder); err != nil {
		return []*Message{}, err
	}

	if im.mboxStatus.Messages == 0 {
		return []*Message{}, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), im.mboxStatus.Messages)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchUid, section.FetchItem()}

	imsg, done := make(chan *imap.Message), make(chan error)
	go func() {
		done <- im.client.Fetch(seqset, items, imsg)
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
		return []*Message{}, fmt.Errorf("%w: %v", ErrIMAPServerProblem, err)
	}

	// for reasons yet unknown, with some imap providers (i.e. Fastmail) the code
	// above sometimes returns the same message twice, but with a different uid.
	dedupMessages := []*Message{}
	for _, m := range messages {
		var isDupe bool
		for _, dm := range dedupMessages {
			if m.Equal(dm) {
				isDupe = true
				break
			}
		}
		if !isDupe {
			dedupMessages = append(dedupMessages, m)
		}
	}

	return dedupMessages, nil
}

func (im *IMAP) Add(folder, subject, body string) error {
	if err := im.Connect(); err != nil {
		return err
	}
	defer im.Close()

	msgStr := fmt.Sprintf("From: %s\r\nDate: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		"gte <gte@ewintr.nl>",
		time.Now().Format(time.RFC822Z),
		subject,
		body,
	)
	msg := NewIMAPBody(msgStr)

	folder = fmt.Sprintf("%s%s", im.config.IMAPFolderPrefix, folder)
	if err := im.client.Append(folder, nil, time.Time{}, imap.Literal(msg)); err != nil {
		return fmt.Errorf("%w: %v", ErrIMAPServerProblem, err)
	}

	return nil
}

func (im *IMAP) Remove(msg *Message) error {
	if msg == nil || !msg.Valid() {
		return ErrInvalidMessage
	}

	if err := im.Connect(); err != nil {
		return err
	}
	defer im.Close()

	if err := im.selectFolder(msg.Folder); err != nil {
		return err
	}

	// set deleted flag
	seqset := new(imap.SeqSet)
	seqset.AddRange(msg.Uid, msg.Uid)
	storeItem := imap.FormatFlagsOp(imap.SetFlags, true)
	err := im.client.UidStore(seqset, storeItem, imap.FormatStringList([]string{imap.DeletedFlag}), nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIMAPServerProblem, err)
	}

	// expunge box
	if err := im.client.Expunge(nil); err != nil {
		return fmt.Errorf("%w: %v", ErrIMAPServerProblem, err)
	}

	return nil
}
