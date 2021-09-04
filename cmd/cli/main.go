package main

import (
	"fmt"
	"os"

	"git.ewintr.nl/gte/cmd/cli/command"
	"git.ewintr.nl/gte/internal/configuration"
)

func main() {
	configPath := "./gte.conf"
	if os.Getenv("GTE_CONFIG") != "" {
		configPath = os.Getenv("GTE_CONFIG")
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err, "could not open config file")
		os.Exit(1)
	}
	config := configuration.New(configFile)

	cmd, err := command.Parse(os.Args[1:], config)
	if err != nil {
		fmt.Println(err, "could not initialize command")
		os.Exit(1)
	}
	fmt.Printf("%s", cmd.Do())
}
