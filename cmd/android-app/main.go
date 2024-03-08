package main

import (
	"code.ewintr.nl/gte/cmd/android-app/component"
	"code.ewintr.nl/gte/cmd/android-app/runner"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
)

func main() {
	fyneApp := app.NewWithID("nl.ewintr.gte")
	fyneWindow := fyneApp.NewWindow("gte - getting things email")
	conf := component.NewConfigurationFromPreferences(fyneApp.Preferences())
	conf.Load()
	logger := component.NewLogger()

	tasksURI, err := storage.Child(fyneApp.Storage().RootURI(), "tasks.json")
	if err != nil {
		logger.Log(err.Error())
	}

	r := runner.NewRunner(conf, tasksURI, logger)
	rootContainer := r.Init()
	fyneWindow.SetContent(rootContainer)
	go r.Run()

	fyneWindow.ShowAndRun()
}
