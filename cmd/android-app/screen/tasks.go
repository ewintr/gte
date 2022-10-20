package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SyncTasksRequest struct{}

type MarkTaskDoneRequest struct {
	ID string
}

type Tasks struct {
	status       binding.String
	tasks        []Task
	taskLabels   binding.StringList
	selectedTask string
	out          chan interface{}
}

func NewTasks(out chan interface{}) *Tasks {
	return &Tasks{
		status:     binding.NewString(),
		tasks:      []Task{},
		taskLabels: binding.NewStringList(),
		out:        out,
	}
}

func (t *Tasks) Refresh(state State) {
	t.status.Set(state.Status)
	t.tasks = state.Tasks
	tls := []string{}
	for _, t := range t.tasks {
		tls = append(tls, t.Action)
	}
	t.taskLabels.Set(tls)
}

func (t *Tasks) Content() fyne.CanvasObject {
	statusLabel := widget.NewLabelWithData(t.status)
	refreshButton := widget.NewButton("refresh", func() {
		t.out <- SyncTasksRequest{}
	})
	doneButton := widget.NewButton("done", func() {
		t.markDone()
	})
	list := widget.NewListWithData(
		t.taskLabels,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	list.OnSelected = t.selectItem

	return container.NewBorder(
		container.NewHBox(refreshButton, statusLabel),
		doneButton,
		nil,
		nil,
		list,
	)
}

func (t *Tasks) selectItem(lid widget.ListItemID) {
	id := int(lid)
	if id < 0 || id >= len(t.tasks) {
		return
	}

	t.selectedTask = t.tasks[id].ID
}

func (t *Tasks) markDone() {
	if t.selectedTask == "" {
		return
	}
	t.out <- MarkTaskDoneRequest{
		ID: t.selectedTask,
	}
}
