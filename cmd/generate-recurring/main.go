package main

import (
	"log"
	"os"
	"strconv"

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
	daysAhead, err := strconv.Atoi(os.Getenv("GTE_DAYS_AHEAD"))
	if err != nil {
		daysAhead = 0
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
	rDate := task.Today.AddDays(daysAhead)
	for _, t := range tasks {
		if t.RecursOn(rDate) {
			subject, body, err := t.CreateDueMessage(rDate)
			if err != nil {
				log.Fatal(err)
			}
			if err := mailStore.Add(task.FOLDER_PLANNED, subject, body); err != nil {
				log.Fatal(err)
			}
		}

	}

}
