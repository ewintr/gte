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
	selectedTask Task
	list         *widget.List
	commands     chan interface{}
	show         chan ShowRequest
	root         *fyne.Container
}

func NewTasks(commands chan interface{}, show chan ShowRequest) *Tasks {
	tasks := &Tasks{
		tasks:        []Task{},
		taskLabels:   binding.NewStringList(),
		commands:     commands,
		show:         show,
		selectedTask: Task{},
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
	if t.selectedTask.ID == "" {
		t.list.UnselectAll()
	}
}

func (t *Tasks) Init() {
	newButton := widget.NewButton("new", func() {
		t.show <- ShowRequest{Screen: "new"}
	})
	updateButton := widget.NewButton("update", func() {
		t.updateTask()
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
		container.NewVBox(updateButton, doneButton),
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

func (t *Tasks) Show(_ Task) {
	t.root.Show()
}

func (t *Tasks) selectItem(lid widget.ListItemID) {
	id := int(lid)
	if id < 0 || id >= len(t.tasks) {
		return
	}

	t.selectedTask = t.tasks[id]
}

func (t *Tasks) markDone() {
	if t.selectedTask.ID == "" {
		return
	}
	t.commands <- MarkTaskDoneRequest{
		ID: t.selectedTask.ID,
	}
	t.selectedTask = Task{}
}

func (t *Tasks) updateTask() {
	if t.selectedTask.ID == "" {
		return
	}

	t.show <- ShowRequest{
		Screen: "update",
		Task:   t.selectedTask,
	}
}
