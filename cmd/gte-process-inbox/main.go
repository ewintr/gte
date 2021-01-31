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
	tasks, err := taskRepo.FindAll(task.FOLDER_INBOX)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		if t.Dirty {
			if err := taskRepo.Update(t); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := taskRepo.CleanUp(); err != nil {
		log.Fatal(err)
	}
}
