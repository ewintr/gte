package runner

import (
	"fmt"
	"time"

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
	r.logger.Log("app started")

	r.requests <- screen.SyncTasksRequest{}
}

func (r *Runner) Run() {
	go r.refresher()
	go r.processRequest()
	go r.backgroundSync()
	r.fyneWindow.ShowAndRun()
}

func (r *Runner) processRequest() {
	for req := range r.requests {
		switch v := req.(type) {
		case screen.SaveConfigRequest:
			r.status = "saving..."
			r.refresh <- true
			for k, val := range v.Fields {
				r.conf.Set(k, val)
			}
			r.logger.Log("new config saved")
			r.status = "ready"
			r.refresh <- true
		case screen.SyncTasksRequest:
			r.status = "syncing..."
			r.refresh <- true
			countDisp, countFetch, err := r.tasks.Sync()
			if err != nil {
				r.logger.Log(err.Error())
			}
			if countDisp > 0 || countFetch > 0 {
				r.logger.Log(fmt.Sprintf("task sync: dispatched: %d, fetched: %d", countDisp, countFetch))
			}
			r.status = "ready"
			r.refresh <- true
		case screen.MarkTaskDoneRequest:
			if err := r.tasks.MarkDone(v.ID); err != nil {
				r.logger.Log(err.Error())
			}
			r.logger.Log(fmt.Sprintf("marked task %q done", v.ID))
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
		sTasks := []screen.Task{}
		for id, action := range tasks {
			sTasks = append(sTasks, screen.Task{
				ID:     id,
				Action: action,
			})
		}

		state := screen.State{
			Status: r.status,
			Tasks:  sTasks,
			Config: r.conf.Fields(),
			Logs:   r.logger.Lines(),
		}
		for _, s := range r.screens {
			s.Refresh(state)
		}
	}
}

func (r *Runner) backgroundSync() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		r.requests <- screen.SyncTasksRequest{}
	}
}
