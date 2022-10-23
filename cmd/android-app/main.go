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

	r := runner.NewRunner(conf, logger)
	tabs := r.Init()
	w.SetContent(tabs)
	go r.Run()

	w.ShowAndRun()
}
