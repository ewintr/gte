package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SyncTasksRequest struct{}

type Tasks struct {
	status binding.String
	tasks  binding.StringList
	out    chan interface{}
}

func NewTasks(out chan interface{}) *Tasks {
	return &Tasks{
		status: binding.NewString(),
		tasks:  binding.NewStringList(),
		out:    out,
	}
}

func (t *Tasks) Refresh(state State) {
	t.status.Set(state.Status)
	t.tasks.Set(state.Tasks)
}

func (t *Tasks) Content() fyne.CanvasObject {
	statusLabel := widget.NewLabelWithData(t.status)
	refreshButton := widget.NewButton("refresh", func() {
		t.out <- SyncTasksRequest{}
	})
	list := widget.NewListWithData(
		t.tasks,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	return container.NewBorder(
		container.NewHBox(refreshButton, statusLabel),
		nil,
		nil,
		nil,
		list,
	)
}
