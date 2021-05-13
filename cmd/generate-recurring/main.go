package main

import (
	"os"
	"strconv"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout).WithField("cmd", "generate-recurring")
	IMAPConfig := &mstore.IMAPConfig{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	}
	msgStore := mstore.NewIMAP(IMAPConfig)

	SMTPConfig := &msend.SSLSMTPConfig{
		URL:      os.Getenv("SMTP_URL"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		To:       os.Getenv("SMTP_TO"),
	}
	if !SMTPConfig.Valid() {
		logger.Error("please set SMTP_URL, SMTP_USERNAME, etc environment variables")
		os.Exit(1)
	}
	mailSend := msend.NewSSLSMTP(SMTPConfig)

	daysAhead, err := strconv.Atoi(os.Getenv("GTE_DAYS_AHEAD"))
	if err != nil {
		daysAhead = 0
	}

	taskRepo := task.NewRepository(msgStore)
	taskDisp := task.NewDispatcher(mailSend)

	recur := process.NewRecur(taskRepo, taskDisp, daysAhead)
	result, err := recur.Process()
	if err != nil {
		logger.WithErr(err).Error("unable to process recurring")
		os.Exit(1)
	}

	logger.WithField("count", result.Count).Info("finished generating recurring tasks")
}
