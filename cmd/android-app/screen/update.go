package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type UpdateTask struct {
	field    *FormField
	action   binding.String
	taskID   string
	commands chan interface{}
	show     chan ShowRequest
	root     *fyne.Container
}

func NewUpdateTask(commands chan interface{}, show chan ShowRequest) *UpdateTask {
	newUpdate := &UpdateTask{
		field:    NewFormField("new due", "due"),
		action:   binding.NewString(),
		commands: commands,
		show:     show,
	}
	newUpdate.Init()

	return newUpdate
}

func (ut *UpdateTask) Init() {
	actionLabel := widget.NewLabel("")
	actionLabel.Bind(ut.action)
	updateForm := widget.NewForm()
	dueEntry := widget.NewEntry()
	dueEntry.Bind(ut.field.Value)
	updateForm.Append(ut.field.Label, dueEntry)

	updateForm.SubmitText = "save"
	updateForm.OnSubmit = ut.Save
	updateForm.CancelText = "cancel"
	updateForm.OnCancel = ut.Cancel
	updateForm.Enable()
	ut.clearForm()

	ut.root = container.NewBorder(
		actionLabel,
		nil,
		nil,
		nil,
		updateForm,
	)
}

func (ut *UpdateTask) Save() {
	ut.commands <- UpdateTaskRequest{
		ID:  ut.taskID,
		Due: ut.field.GetValue(),
	}
	ut.show <- ShowRequest{Screen: "tasks"}
}

func (ut *UpdateTask) Cancel() {
	ut.clearForm()
	ut.show <- ShowRequest{Screen: "tasks"}
}

func (ut *UpdateTask) clearForm() {
	ut.field.SetValue("")
	ut.action.Set("")
	ut.taskID = ""
}

func (ut *UpdateTask) Refresh(_ State) {}

func (ut *UpdateTask) Content() *fyne.Container {
	return ut.root
}

func (ut *UpdateTask) Hide() {
	ut.root.Hide()
}

func (ut *UpdateTask) Show(task Task) {
	ut.taskID = task.ID
	ut.action.Set(task.Action)
	ut.root.Show()
}
