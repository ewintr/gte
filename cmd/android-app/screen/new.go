package screen

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var newLock sync.Mutex

type NewTask struct {
	fields   []*FormField
	commands chan interface{}
	show     chan string
	root     *fyne.Container
}

func NewNewTask(commands chan interface{}, show chan string) *NewTask {
	fields := []*FormField{}
	for _, f := range [][2]string{
		{"action", "action"},
		{"project", "project"},
		{"due", "due string"},
		{"recur", "recur string"},
	} {
		fields = append(fields, NewFormField(f[0], f[1]))
	}

	newTask := &NewTask{
		fields:   fields,
		commands: commands,
		show:     show,
	}
	newTask.Init()

	return newTask
}

func (nt *NewTask) Init() {
	taskForm := widget.NewForm()
	for _, f := range nt.fields {
		w := widget.NewEntry()
		w.Bind(f.Value)
		taskForm.Append(f.Label, w)
	}

	taskForm.SubmitText = "save"
	taskForm.OnSubmit = nt.Save
	taskForm.Enable()
	nt.clearForm()

	nt.root = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		taskForm,
	)
}

func (nt *NewTask) Save() {
	req := SaveNewTaskRequest{
		Fields: map[string]string{},
	}
	for _, f := range nt.fields {
		req.Fields[f.Key] = f.GetValue()
	}
	nt.commands <- req
	nt.show <- "tasks"

	nt.clearForm()
}

func (nt *NewTask) clearForm() {
	for _, f := range nt.fields {
		f.SetValue("")
	}
}

func (nt *NewTask) Refresh(_ State) {}

func (nt *NewTask) Content() *fyne.Container {
	return nt.root
}

func (nt *NewTask) Hide() {
	nt.root.Hide()
}

func (nt *NewTask) Show() {
	nt.root.Show()
}
