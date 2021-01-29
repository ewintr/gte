package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~ewintr/gte/internal/task"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

func main() {
	config := &mstore.EmailConfiguration{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	}
	if !config.Valid() {
		log.Fatal("please set MAIL_USER, MAIL_PASSWORD, etc environment variables")
	}

	mailStore, err := mstore.EmailConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mailStore.Disconnect()

	taskRepo := task.NewRepository(mailStore)
	tasks, err := taskRepo.FindAll("INBOX")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		fmt.Printf("processing: %s... ", t.Action)

		if t.Dirty() {
			if err := taskRepo.Update(t); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("updated.")
		}
		fmt.Printf("\n")
	}

	/*
			folders, err := mailStore.FolderNames()
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range folders {
				fmt.Println(f)
			}

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

				body := NewBody(`From: todo <process@erikwinter.nl>
			Subject: the subject

			And here comes the body`)

				if err := mailStore.Append("INBOX", imap.Literal(body)); err != nil {
					log.Fatal(err)
				}
	*/
}
