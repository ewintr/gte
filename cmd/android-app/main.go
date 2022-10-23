package main

import (
	"ewintr.nl/gte/cmd/android-app/component"
	"ewintr.nl/gte/cmd/android-app/runner"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
)

func main() {
	fyneApp := app.NewWithID("nl.ewintr.gte")
	w := fyneApp.NewWindow("gte - getting things email")
	conf := component.NewConfigurationFromPreferences(fyneApp.Preferences())
	conf.Load()
	logger := component.NewLogger()

	tasksURI, err := storage.Child(fyneApp.Storage().RootURI(), "tasks.json")
	if err != nil {
		logger.Log(err.Error())
	}

	r := runner.NewRunner(conf, tasksURI, logger)
	tabs := r.Init()
	w.SetContent(tabs)
	go r.Run()

	w.ShowAndRun()
}
