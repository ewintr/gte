package screen

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Tasks struct {
	tasks        []Task
	taskLabels   binding.StringList
	selectedTask string
	list         *widget.List
	commands     chan interface{}
	show         chan string
	root         *fyne.Container
}

func NewTasks(commands chan interface{}, show chan string) *Tasks {
	tasks := &Tasks{
		tasks:      []Task{},
		taskLabels: binding.NewStringList(),
		commands:   commands,
		show:       show,
	}
	tasks.Init()

	return tasks
}

func (t *Tasks) Refresh(state State) {
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

func (t *Tasks) Init() {
	newButton := widget.NewButton("new", func() {
		t.show <- "new"
	})
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

	t.root = container.NewBorder(
		newButton,
		doneButton,
		nil,
		nil,
		t.list,
	)
}

func (t *Tasks) Content() *fyne.Container {
	return t.root
}

func (t *Tasks) Hide() {
	t.root.Hide()
}

func (t *Tasks) Show() {
	t.root.Show()
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
	t.commands <- MarkTaskDoneRequest{
		ID: t.selectedTask,
	}
	t.selectedTask = ""
}
