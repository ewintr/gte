package main

import (
	"os"
	"os/signal"
	"time"

	"ewintr.nl/go-kit/log"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
	"ewintr.nl/gte/pkg/msend"
	"ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout)

	config := configuration.NewFromEnvironment()
	cfgCopy := *config
	cfgCopy.IMAPPassword = "***"
	cfgCopy.SMTPPassword = "***"
	logger.WithField("config", cfgCopy).Info("started")

	msgStore := mstore.NewIMAP(config.IMAP())
	mailSend := msend.NewSSLSMTP(config.SMTP())
	repo := storage.NewRemoteRepository(msgStore)
	disp := storage.NewDispatcher(mailSend)

	inboxProc := process.NewInbox(repo)
	recurProc := process.NewRecur(repo, disp, config.DaysAhead)

	go Run(inboxProc, recurProc, logger)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)
	<-done
	logger.Info("stopped")
}

func Run(inboxProc *process.Inbox, recurProc *process.Recur, logger log.Logger) {
	logger = logger.WithField("func", "run")
	inboxTicker := time.NewTicker(30 * time.Second)
	recurTicker := time.NewTicker(time.Hour)
	oldToday := task.Today()

	for {
		select {
		case <-inboxTicker.C:
			result, err := inboxProc.Process()
			if err != nil {
				logger.WithErr(err).Error("failed processing inbox")

				continue
			}
			if result.Count > 0 {
				logger.WithField("result", result).Info("finished processing inbox")
			}
		case <-recurTicker.C:
			if oldToday.Equal(task.Today()) {

				continue
			}

			oldToday = task.Today()
			result, err := recurProc.Process()
			if err != nil {
				logger.WithErr(err).Error("failed generating recurring tasks")

				continue
			}
			if result.Count > 0 {
				logger.WithField("result", result).Info("finished generating recurring tasks")
			}
		}
	}
}
