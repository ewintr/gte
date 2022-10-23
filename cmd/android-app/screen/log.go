package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Log struct {
	lines binding.StringList
}

func NewLog() *Log {
	return &Log{
		lines: binding.NewStringList(),
	}
}

func (l *Log) Refresh(state State) {
	l.lines.Set(state.Logs)
}

func (l *Log) Content() fyne.CanvasObject {
	list := widget.NewListWithData(
		l.lines,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	return list
}
