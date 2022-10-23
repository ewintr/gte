package screen

import (
	"fmt"
	"sort"

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
	list         *widget.List
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
	t.status.Set(fmt.Sprintf("> %s", state.Status))
	t.tasks = state.Tasks
	sort.Slice(t.tasks, func(i, j int) bool {
		return t.tasks[i].Action < t.tasks[j].Action
	})
	tls := []string{}
	for _, t := range t.tasks {
		tls = append(tls, t.Action)
	}
	t.taskLabels.Set(tls)
	if t.selectedTask == "" {
		t.list.UnselectAll()
	}
}

func (t *Tasks) Content() fyne.CanvasObject {
	statusLabel := widget.NewLabel("> init...")
	statusLabel.Bind(t.status)
	statusLabel.TextStyle.Italic = true
	doneButton := widget.NewButton("done", func() {
		t.markDone()
	})
	t.list = widget.NewListWithData(
		t.taskLabels,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	t.list.OnSelected = t.selectItem

	return container.NewBorder(
		container.NewHBox(statusLabel),
		doneButton,
		nil,
		nil,
		t.list,
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
	t.selectedTask = ""
}
