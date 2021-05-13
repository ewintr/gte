package main

import (
	"os"
	"os/signal"
	"time"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout)
	logger.Info("started")

	msgStore := mstore.NewIMAP(&mstore.IMAPConfig{
		IMAPURL:      os.Getenv("IMAP_URL"),
		IMAPUsername: os.Getenv("IMAP_USERNAME"),
		IMAPPassword: os.Getenv("IMAP_PASSWORD"),
	})
	msgSender := msend.NewSSLSMTP(&msend.SSLSMTPConfig{
		URL:      os.Getenv("SMTP_URL"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		To:       os.Getenv("SMTP_TO"),
	})
	repo := task.NewRepository(msgStore)
	disp := task.NewDispatcher(msgSender)

	inboxProc := process.NewInbox(repo)
	recurProc := process.NewRecur(repo, disp, 6)

	go Run(inboxProc, recurProc, logger)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)
	<-done
	logger.Info("stopped")
}

func Run(inboxProc *process.Inbox, recurProc *process.Recur, logger log.Logger) {
	logger = logger.WithField("func", "run")
	inboxTicker := time.NewTicker(10 * time.Second)
	recurTicker := time.NewTicker(time.Hour)
	oldToday := task.Today

	for {
		select {
		case <-inboxTicker.C:
			result, err := inboxProc.Process()
			if err != nil {
				logger.WithErr(err).Error("failed processing inbox")

				continue
			}
			logger.WithField("count", result.Count).Info("finished processing inbox")
		case <-recurTicker.C:
			year, month, day := time.Now().Date()
			newToday := task.NewDate(year, int(month), day)
			if oldToday.Equal(newToday) {

				continue
			}

			oldToday = newToday
			result, err := recurProc.Process()
			if err != nil {
				logger.WithErr(err).Error("failed generating recurring tasks")

				continue
			}
			logger.WithField("count", result.Count).Info("finished generating recurring tasks")
		}
	}
}
