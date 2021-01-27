package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

func main() {
	iPort, _ := strconv.Atoi(os.Getenv("IMAP_PORT"))
	//sPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	config := &mstore.EMailStoreConfiguration{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPPort:     iPort,
		IMAPUsername: os.Getenv("MAIL_USER"),
		IMAPPassword: os.Getenv("MAIL_PASSWORD"),
	}
	if !config.Valid() {
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

		// List mailboxes
		mailboxes := make(chan *imap.MailboxInfo, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.List("", "*", mailboxes)
		}()

		log.Println("Mailboxes:")
		for m := range mailboxes {
			log.Println("* " + m.Name)
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}

		// Select INBOX
		mbox, err := c.Select("INBOX", false)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Flags for INBOX:", mbox.Flags)

		// Get the last 4 messages
		from := uint32(1)
		to := mbox.Messages
		if mbox.Messages > 3 {
			// We're using unsigned integers here, only subtract if the result is > 0
			from = mbox.Messages - 3
		}
		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)

		messages := make(chan *imap.Message, 10)
		done = make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		log.Println("Last 4 messages:")
		for msg := range messages {
			log.Println("* " + msg.Envelope.Subject)
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}

		log.Println("Done!")
	*/
}
