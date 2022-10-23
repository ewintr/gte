package runner

import (
	"fmt"
	"sync"
	"time"

	"ewintr.nl/gte/cmd/android-app/component"
	"ewintr.nl/gte/cmd/android-app/screen"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var runnerLock = sync.Mutex{}

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

func NewRunner(conf *component.Configuration, logger *component.Logger) *Runner {
	return &Runner{
		status:   "init",
		conf:     conf,
		logger:   logger,
		requests: make(chan interface{}),
		refresh:  make(chan bool),
	}
}

func (r *Runner) Init() fyne.CanvasObject {

	logScreen := screen.NewLog()
	logTab := container.NewTabItem("log", logScreen.Content())
	r.screens = append(r.screens, logScreen)

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

	return tabs
}

func (r *Runner) Run() {
	go r.refresher()
	go r.processRequest()
	r.backgroundSync()
}

func (r *Runner) processRequest() {
	for req := range r.requests {
		r.logger.Log(fmt.Sprintf("processing request %T", req))
		switch v := req.(type) {
		case screen.SaveConfigRequest:
			r.status = "saving..."
			r.refresh <- true
			for k, val := range v.Fields {
				r.conf.Set(k, val)
			}
			r.logger.Log("new config saved")
			r.status = "ready"
		case screen.SyncTasksRequest:
			r.logger.Log("starting sync request")
			r.status = "syncing..."
			r.refresh <- true
			countDisp, countFetch, err := r.tasks.Sync(r.logger)
			if err != nil {
				r.logger.Log(err.Error())
			}
			//if countDisp > 0 || countFetch > 0 {
			r.logger.Log(fmt.Sprintf("task sync: dispatched: %d, fetched: %d", countDisp, countFetch))
			//}
			r.status = "ready"
			r.logger.Log("sync request done")

		case screen.MarkTaskDoneRequest:
			if err := r.tasks.MarkDone(v.ID); err != nil {
				r.logger.Log(err.Error())
			}
			r.logger.Log(fmt.Sprintf("marked task %q done", v.ID))
		default:
			r.logger.Log("request unknown")
		}
		r.refresh <- true
		r.logger.Log("processing request done")
	}

}

func (r *Runner) refresher() {
	for <-r.refresh {
		r.logger.Log("start refresh")
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
	ticker := time.NewTicker(15 * time.Second)
	for {
		r.requests <- screen.SyncTasksRequest{}
		<-ticker.C
	}
}
