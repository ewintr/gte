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
	if len(l.lines) > 50 {
		l.lines = l.lines[1:]
	}
}

func (l *Logger) Lines() []string {
	return l.lines
}
