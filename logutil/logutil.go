package logutil

import (
	"fmt"
	"strings"
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

const indent = "   "

type Writer struct {
	Source      string
	TaskID      string
	DefaultType int
	depth       int
}

func (l *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	msg := string(p[:n])
	msgs := strings.Split(msg, "\n")
	for _, m := range msgs {
		l.WriteType(m, l.DefaultType)
	}
	return n, nil
}

func (l *Writer) WriteType(msg string, logType int) {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now().Format("2006-01-02T15:04:05")
	indents := ""
	for i := 0; i < l.depth; i++ {
		indents += indent
	}

	fmt.Printf("[%s][%s][%d] %s\t| %s%s\n", l.Source, l.TaskID, logType, now, indents, msg)
}

func (l *Writer) Indent() {
	l.depth += 1
}

func (l *Writer) Outdent() {
	l.depth -= 1
	if l.depth < 0 {
		l.depth = 0
	}
}

func (l *Writer) CreateChild() *Writer {
	return &Writer{
		Source: l.Source,
		TaskID: l.TaskID,
		depth:  l.depth + 1,
	}
}

func (l *Writer) CreateWriter(logType int) (w *Writer) {
	return &Writer{
		Source:      l.Source,
		TaskID:      l.TaskID,
		DefaultType: logType,
		depth:       l.depth,
	}
}

func NewWriter(source string, taskID string) *Writer {
	return &Writer{
		Source: source,
		TaskID: taskID,
	}
}
