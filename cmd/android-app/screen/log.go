package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Log struct {
	lines binding.StringList
	root  *fyne.Container
}

func NewLog() *Log {
	logs := &Log{
		lines: binding.NewStringList(),
	}
	logs.Init()

	return logs
}

func (l *Log) Refresh(state State) {
	l.lines.Set(state.Logs)
}

func (l *Log) Init() {
	list := widget.NewListWithData(
		l.lines,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	l.root = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		list,
	)
}

func (l *Log) Content() *fyne.Container {
	return l.root
}

func (l *Log) Hide() {
	l.root.Hide()
}

func (l *Log) Show(_ Task) {
	l.root.Show()
}
