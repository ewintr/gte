package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

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

	messages, err := mailStore.Inbox()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range messages {
		fmt.Println(m.Subject)
	}

}
