package main

import (
	"os"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout).WithField("cmd", "process-inbox")
	config := &mstore.IMAPConfig{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	}
	msgStore := mstore.NewIMAP(config)

	inboxProcessor := process.NewInbox(task.NewRepository(msgStore))
	result, err := inboxProcessor.Process()
	if err != nil {
		logger.WithErr(err).Error("unable to process inbox")
		os.Exit(1)
	}
	logger.WithField("count", result.Count).Info("finished processing inbox")
}
