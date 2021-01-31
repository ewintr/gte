package main

import (
	"log"
	"os"

	"git.sr.ht/~ewintr/gte/internal/task"
	"git.sr.ht/~ewintr/gte/pkg/mstore"
)

func main() {
	config := &mstore.ImapConfiguration{
		ImapUrl:      os.Getenv("IMAP_URL"),
		ImapUsername: os.Getenv("IMAP_USERNAME"),
		ImapPassword: os.Getenv("IMAP_PASSWORD"),
	}
	if !config.Valid() {
		log.Fatal("please set IMAP_USER, IMAP_PASSWORD, etc environment variables")
	}

	mailStore, err := mstore.ImapConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mailStore.Disconnect()

	taskRepo := task.NewRepository(mailStore)
	tasks, err := taskRepo.FindAll(task.FOLDER_RECURRING)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		if t.RecursToday() {
			subject, body, err := t.CreateNextMessage(task.Today)
			if err != nil {
				log.Fatal(err)
			}
			if err := mailStore.Add(task.FOLDER_PLANNED, subject, body); err != nil {
				log.Fatal(err)
			}
		}

	}

}
