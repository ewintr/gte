package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout)

	configPath := flag.String("c", "~/.config/gte/gte.conf", "path to configuration file")
	daysAhead := flag.Int("d", 6, "generate for this amount of days from now")
	flag.Parse()

	logger.With(log.Fields{
		"config":    *configPath,
		"daysAhead": *daysAhead,
	}).Info("started")

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.WithErr(err).Error("could not open config file")
		os.Exit(1)
	}
	config := configuration.New(configFile)

	msgStore := mstore.NewIMAP(config.IMAP())
	mailSend := msend.NewSSLSMTP(config.SMTP())
	repo := task.NewRepository(msgStore)
	disp := task.NewDispatcher(mailSend)

	inboxProc := process.NewInbox(repo)
	recurProc := process.NewRecur(repo, disp, *daysAhead)

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
