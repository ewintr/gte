package main

import (
	"flag"
	"os"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/msend"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout).WithField("cmd", "generate-recurring")

	configPath := flag.String("c", "~/.config/gte/gte.conf", "path to configuration file")
	daysAhead := flag.Int("d", 6, "generate for this amount of days from now")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.WithErr(err).Error("could not open config file")
		os.Exit(1)
	}
	config := configuration.New(configFile)

	msgStore := mstore.NewIMAP(config.IMAP())
	mailSend := msend.NewSSLSMTP(config.SMTP())
	taskRepo := task.NewRemoteRepository(msgStore)
	taskDisp := task.NewDispatcher(mailSend)
	recur := process.NewRecur(taskRepo, taskDisp, *daysAhead)

	result, err := recur.Process()
	if err != nil {
		logger.WithErr(err).Error("unable to process recurring")
		os.Exit(1)
	}

	logger.WithField("result", result).Info("finished generating recurring tasks")
}
