package logutil

import (
	"fmt"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

const (
	Verbose int = iota
	Info
	Warning
	Error
)

type Logger struct {
	TaskID string
	depth  int
}

func (l *Logger) Write(msg string, logType int) {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now().Format("2006-01-02T15:04:05")
	indent := ""
	for i := 0; i < l.depth; i++ {
		indent += "    "
	}
	fmt.Printf("[%s][%d] %s\t| %s%s\n", l.TaskID, logType, now, indent, msg)
}

func (l *Logger) CreateChild() *Logger {
	return &Logger{TaskID: l.TaskID, depth: l.depth + 1}
}

func NewLogger(taskID string) *Logger {
	return &Logger{TaskID: taskID}
}
