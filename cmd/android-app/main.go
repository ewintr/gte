package main

import (
	"ewintr.nl/gte/cmd/android-app/component"
	"ewintr.nl/gte/cmd/android-app/runner"
	"fyne.io/fyne/v2/app"
)

func main() {
	fyneApp := app.NewWithID("nl.ewintr.gte")
	w := fyneApp.NewWindow("gte - getting things email")
	conf := component.NewConfigurationFromPreferences(fyneApp.Preferences())
	conf.Load()
	logger := component.NewLogger()

	runner := runner.NewRunner(conf, logger)
	tabs := runner.Init()
	w.SetContent(tabs)
	go runner.Run()

	w.ShowAndRun()
}
