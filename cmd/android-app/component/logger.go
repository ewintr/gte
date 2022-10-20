package component

import (
	"fmt"
	"time"
)

type Logger struct {
	lines []string
}

func NewLogger() *Logger {
	return &Logger{
		lines: []string{},
	}
}

func (l *Logger) Log(line string) {
	l.lines = append(l.lines, fmt.Sprintf("%s: %s", time.Now().Format("15:04:05"), line))
}

func (l *Logger) Lines() []string {
	if len(l.lines) == 0 {
		return []string{}
	}

	last := len(l.lines) - 1
	first := last - 50
	if first < 0 {
		first = 0
	}
	reverse := []string{}
	for i := last; i >= first; i-- {
		reverse = append(reverse, l.lines[i])
	}
	return reverse
}
