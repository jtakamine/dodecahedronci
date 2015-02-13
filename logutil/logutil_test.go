package logutil

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	id := "asdlfjalsdjflajdlkajlsjdli12342;kj;23l4ll23"
	l := NewLogger(id)
	if l.TaskID != id {
		t.Errorf("NewLogger(\"%s\") returned Logger with TaskID=\"%s\"")
	}
}

func TestWrite(t *testing.T) {
	l := NewLogger("12345asdf")
	l.Write("top level verbose message", Verbose)
	l.Write("top level info message", Info)
	l.Write("top level warning message", Warning)
	l.Write("top level error message", Error)

	l2 := l.CreateChild()
	l2.Write("child (depth=1) info message", Info)
	l2.Write("child (depth=1) info message", Info)

	l3 := l2.CreateChild()
	l3.Write("child (depth=2) warning message", Warning)

	l2.Write("child (depth=1) info message", Info)

	l.Write("(last) top level info message", Info)
}
