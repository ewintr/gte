package runner

import (
	"fmt"

	"ewintr.nl/gte/cmd/android-app/component"
	"ewintr.nl/gte/cmd/android-app/screen"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type Runner struct {
	fyneApp    fyne.App
	fyneWindow fyne.Window
	conf       *component.Configuration
	logger     *component.Logger
	tasks      *component.Tasks
	screens    []screen.Screen
	status     string
	requests   chan interface{}
	refresh    chan bool
}

func NewRunner() *Runner {
	return &Runner{
		requests: make(chan interface{}, 3),
		refresh:  make(chan bool),
	}
}

func (r *Runner) Init() {
	fyneApp := app.NewWithID("nl.ewintr.gte")
	w := fyneApp.NewWindow("gte - getting things email")

	r.logger = component.NewLogger()
	logScreen := screen.NewLog()
	logTab := container.NewTabItem("log", logScreen.Content())
	r.screens = append(r.screens, logScreen)

	r.conf = component.NewConfigurationFromPreferences(fyneApp.Preferences())
	r.conf.Load()
	configScreen := screen.NewConfig(r.requests)
	configTab := container.NewTabItem("config", configScreen.Content())
	r.screens = append(r.screens, configScreen)

	r.logger.Log("initializing tasks...")
	tasks, err := component.NewTasks(r.conf)
	if err != nil {
		r.logger.Log(err.Error())
	}
	r.tasks = tasks
	taskScreen := screen.NewTasks(r.requests)
	taskTab := container.NewTabItem("tasks", taskScreen.Content())
	r.screens = append(r.screens, taskScreen)
	tabs := container.NewAppTabs(taskTab, configTab, logTab)

	w.SetContent(tabs)

	r.fyneApp = fyneApp
	r.fyneWindow = w

	r.requests <- screen.SyncTasksRequest{}
}

func (r *Runner) Run() {
	go r.refresher()
	go r.processRequest()
	r.fyneWindow.ShowAndRun()
}

func (r *Runner) processRequest() {
	for req := range r.requests {
		r.logger.Log(fmt.Sprintf("request %T: %s", req, req))

		switch v := req.(type) {
		case screen.SaveConfigRequest:
			for k, val := range v.Fields {
				r.conf.Set(k, val)
			}
		case screen.SyncTasksRequest:
			r.status = "syncing..."
			r.refresh <- true
			count, err := r.tasks.Sync()
			if err != nil {
				r.logger.Log(err.Error())
			}
			r.logger.Log(fmt.Sprintf("fetched: %d", count))
			r.status = "synced"
			r.refresh <- true
		default:
			r.logger.Log("request unknown")
		}
	}
}

func (r *Runner) refresher() {
	for range r.refresh {
		tasks, err := r.tasks.Today()
		if err != nil {
			r.logger.Log(err.Error())
		}
		state := screen.State{
			Status: r.status,
			Tasks:  tasks,
			Config: r.conf.Fields(),
			Logs:   r.logger.Lines(),
		}

		for _, s := range r.screens {
			s.Refresh(state)
		}
	}
}
