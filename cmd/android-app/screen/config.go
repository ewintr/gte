package screen

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SaveConfigRequest struct {
	Fields map[string]string
}

var confLock sync.Mutex

type FormField struct {
	Label string
	Key   string
	Value binding.String
}

func NewFormField(key, label string) *FormField {
	val := binding.NewString()
	val.Set("...")

	return &FormField{
		Label: label,
		Key:   key,
		Value: val,
	}
}

func (ff *FormField) SetValue(value string) {
	confLock.Lock()
	defer confLock.Unlock()

	ff.Value.Set(value)
}

func (ff *FormField) GetValue() string {
	val, _ := ff.Value.Get()

	return val
}

type Config struct {
	fields []*FormField
	out    chan interface{}
}

func NewConfig(out chan interface{}) *Config {
	fields := []*FormField{}
	for _, f := range [][2]string{
		{"ConfigIMAPURL", "imap url"},
		{"ConfigIMAPUser", "imap user"},
		{"ConfigIMAPPassword", "imap password"},
		{"ConfigIMAPFolderPrefix", "imap folder prefix"},
		{"ConfigSMTPURL", "smtp url"},
		{"ConfigSMTPUser", "smtp user"},
		{"ConfigSMTPPassword", "smtp password"},
		{"ConfigGTEToName", "to name"},
		{"ConfigGTEToAddress", "to address"},
		{"ConfigGTEFromName", "from name"},
		{"ConfigGTEFromAddress", "from address"},
		{"ConfigGTELocalDBPath", "local db path"},
	} {
		fields = append(fields, NewFormField(f[0], f[1]))
	}

	return &Config{
		fields: fields,
		out:    out,
	}
}

func (cf *Config) Save() {
	req := SaveConfigRequest{
		Fields: map[string]string{},
	}
	for _, f := range cf.fields {
		req.Fields[f.Key] = f.GetValue()
	}
	cf.out <- req
}

func (cf *Config) Refresh(state State) {
	for _, f := range cf.fields {
		if v, ok := state.Config[f.Key]; ok {
			f.SetValue(v)
		}
	}
}

func (cf *Config) Content() fyne.CanvasObject {
	configForm := widget.NewForm()
	for _, f := range cf.fields {
		w := widget.NewEntry()
		w.Bind(f.Value)
		configForm.Append(f.Label, w)
	}

	configForm.SubmitText = "save"
	configForm.OnSubmit = cf.Save
	configForm.Enable()

	return configForm
}
