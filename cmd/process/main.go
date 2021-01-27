package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
	"github.com/emersion/go-imap"
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

func main() {
	config := &mstore.EMailStoreConfiguration{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	}
	if !config.Valid() {
		fmt.Printf("conf: %v\n", config)
		log.Fatal("please set MAIL_USER, MAIL_PASSWORD, etc environment variables")
	}
	//fmt.Printf("conf: %+v\n", config)

	mailStore, err := mstore.EMailConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mailStore.Disconnect()

	folders, err := mailStore.Folders()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range folders {
		fmt.Println(f.Name)
	}

	/*
		messages, err := mailStore.Inbox()
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range messages {
			fmt.Println(m.Subject)
		}
	*/

	body := NewBody(`From: todo <process@erikwinter.nl>
Subject: the subject

And here comes the body`)

	if err := mailStore.Append("INBOX", imap.Literal(body)); err != nil {
		log.Fatal(err)
	}
}
