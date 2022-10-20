package screen

import "fyne.io/fyne/v2"

type Task struct {
	ID     string
	Action string
}

type State struct {
	Status string
	Tasks  []Task
	Config map[string]string
	Logs   []string
}

type Screen interface {
	Content() fyne.CanvasObject
	Refresh(state State)
}
