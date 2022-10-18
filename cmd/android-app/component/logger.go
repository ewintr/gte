package component

type Logger struct {
	lines []string
}

func NewLogger() *Logger {
	return &Logger{
		lines: []string{},
	}
}

func (l *Logger) Log(line string) {
	l.lines = append(l.lines, line)
}

func (l *Logger) Lines() []string {
	return l.lines
}
