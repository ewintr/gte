package screen

import (
	"sync"

	"fyne.io/fyne/v2/data/binding"
)

var formLock sync.Mutex

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
	formLock.Lock()
	defer formLock.Unlock()

	ff.Value.Set(value)
}

func (ff *FormField) GetValue() string {
	val, _ := ff.Value.Get()

	return val
}
