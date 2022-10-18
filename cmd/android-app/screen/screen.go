package screen

import "fyne.io/fyne/v2"

type State struct {
	Status string
	Tasks  []string
	Config map[string]string
	Logs   []string
}

type Screen interface {
	Content() fyne.CanvasObject
	Refresh(state State)
}
