package main

import (
	"flag"
	"os"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

func main() {
	logger := log.New(os.Stdout).WithField("cmd", "process-inbox")

	configPath := flag.String("c", "~/.config/gte/gte.conf", "path to configuration file")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.WithErr(err).Error("could not open config file")
		os.Exit(1)
	}
	config := configuration.New(configFile)
	msgStore := mstore.NewIMAP(config.IMAP())
	inboxProcessor := process.NewInbox(task.NewRepository(msgStore))

	result, err := inboxProcessor.Process()
	if err != nil {
		logger.WithErr(err).Error("unable to process inbox")
		os.Exit(1)
	}
	logger.WithField("result", result).Info("finished processing inbox")
}
