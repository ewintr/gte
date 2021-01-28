package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
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
	config := &mstore.EmailConfiguration{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	}
	if !config.Valid() {
		fmt.Printf("conf: %v\n", config)
		log.Fatal("please set MAIL_USER, MAIL_PASSWORD, etc environment variables")
	}
	//fmt.Printf("conf: %+v\n", config)

	mailStore, err := mstore.EmailConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mailStore.Disconnect()

	/*
		folders, err := mailStore.FolderNames()
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range folders {
			fmt.Println(f)
		}
	*/

	if err := mailStore.Select("Today"); err != nil {
		log.Fatal(err)
	}

	messages, err := mailStore.Messages()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range messages {
		fmt.Printf("%d: %s\n", m.Uid, m.Subject)
	}
	if len(messages) == 0 {
		log.Fatal("no messages")
		return
	}

	if err := mailStore.Remove(messages[0].Uid); err != nil {
		log.Fatal(err)
	}

	/*
			body := NewBody(`From: todo <process@erikwinter.nl>
		Subject: the subject

		And here comes the body`)

			if err := mailStore.Append("INBOX", imap.Literal(body)); err != nil {
				log.Fatal(err)
			}
	*/
}
