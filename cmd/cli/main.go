package main

import (
	"fmt"
	"os"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/gte/cmd/cli/command"
	"git.ewintr.nl/gte/internal/configuration"
)

func main() {
	loglevel := log.LogLevel("error")
	if os.Getenv("GTE_LOGLEVEL") != "" {
		loglevel = log.LogLevel(os.Getenv("GTE_LOGLEVEL"))
	}
	logger := log.New(os.Stdout).WithField("cmd", "cli")
	logger.SetLogLevel(loglevel)

	configPath := "/home/erik/.config/gte/gte.conf"
	if os.Getenv("GTE_CONFIG") != "" {
		configPath = os.Getenv("GTE_CONFIG")
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		logger.WithErr(err).Error("could not open config file")
		os.Exit(1)
	}
	config := configuration.New(configFile)

	cmd, err := command.Parse(os.Args[1:], config)
	if err != nil {
		logger.WithErr(err).Error("could not initialize command")
		os.Exit(1)
	}
	result, err := cmd.Do()
	if err != nil {
		logger.WithErr(err).Error("could perform command")
		os.Exit(1)
	}
	fmt.Printf("%s\n", result.Message)
}
