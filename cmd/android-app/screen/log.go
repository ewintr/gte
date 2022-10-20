package screen

import (
	"sort"

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
	lines := state.Logs
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] > lines[j]
	})
	l.lines.Set(lines)
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
