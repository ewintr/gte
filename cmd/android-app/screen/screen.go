package screen

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	ID     string
	Action string
}

type ShowRequest struct {
	Screen string
	Task   Task
}

type State struct {
	Status        string
	CurrentScreen string
	Tasks         []Task
	Config        map[string]string
	Logs          []string
}

type Screen interface {
	Content() *fyne.Container
	Refresh(state State)
	Hide()
	Show(Task)
}

type ScreenSet struct {
	current string
	show    chan ShowRequest
	status  binding.String
	menu    *fyne.Container
	screens map[string]Screen
	root    *fyne.Container
}

func NewScreenSet(requests chan interface{}) *ScreenSet {
	status := binding.NewString()
	show := make(chan ShowRequest)

	tasksButton := widget.NewButton("tasks", func() {
		show <- ShowRequest{Screen: "tasks"}
	})
	configButton := widget.NewButton("config", func() {
		show <- ShowRequest{Screen: "config"}
	})
	logsButton := widget.NewButton("logs", func() {
		show <- ShowRequest{Screen: "logs"}
	})
	statusLabel := widget.NewLabel("> init...")
	statusLabel.Bind(status)
	statusLabel.TextStyle.Italic = true
	menu := container.NewHBox(tasksButton, configButton, logsButton, statusLabel)

	screens := map[string]Screen{
		"tasks":  NewTasks(requests, show),
		"logs":   NewLog(),
		"config": NewConfig(requests, show),
		"new":    NewNewTask(requests, show),
		"update": NewUpdateTask(requests, show),
	}

	cs := []fyne.CanvasObject{}
	for _, s := range screens {
		s.Hide()
		cs = append(cs, s.Content())
	}
	screens["tasks"].Show(Task{})

	root := container.NewBorder(menu, nil, nil, nil, cs...)

	return &ScreenSet{
		status:  status,
		current: "tasks",
		show:    show,
		screens: screens,
		root:    root,
	}
}

func (ss *ScreenSet) Run() {
	for s := range ss.show {
		if s.Screen != ss.current {
			ss.screens[ss.current].Hide()
			ss.screens[s.Screen].Show(s.Task)
			ss.current = s.Screen

			ss.root.Refresh()
		}
	}
}

func (ss *ScreenSet) Refresh(state State) {
	ss.status.Set(fmt.Sprintf("> %s", state.Status))
	for _, s := range ss.screens {
		s.Refresh(state)
	}

	ss.root.Refresh()
}

func (ss *ScreenSet) Content() *fyne.Container {
	return ss.root
}
