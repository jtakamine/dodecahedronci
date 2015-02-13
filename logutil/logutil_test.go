package logutil

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	id := "asdlfjalsdjflajdlkajlsjdli12342;kj;23l4ll23"
	w := NewWriter(id)
	if w.TaskID != id {
		t.Errorf("NewLogger(\"%s\") returned Logger with TaskID=\"%s\"")
	}
}

func TestChildWriteType(t *testing.T) {
	w := NewWriter("12345asdf")
	w.WriteType("top level verbose message", Verbose)
	w.WriteType("top level info message", Info)
	w.WriteType("top level warning message", Warning)
	w.WriteType("top level error message", Error)

	w2 := w.CreateChild()
	w2.WriteType("child (depth=1) info message", Info)
	w2.WriteType("child (depth=1) info message", Info)

	w3 := w2.CreateChild()
	w3.WriteType("child (depth=2) warning message", Warning)

	w2.WriteType("child (depth=1) info message", Info)

	w.WriteType("(last) top level info message", Info)
}

func TestIndentOutdent(t *testing.T) {
	w := NewWriter("987663ASDF")
	w.WriteType("non-indented message", Info)
	w.WriteType("non-indented message #2", Info)
	w.Indent()
	w.WriteType("indented once message", Info)
	w.WriteType("indented once message", Info)
	w.Indent()
	w.WriteType("indented twice message", Info)
	w.Outdent()
	w.WriteType("indented once message", Info)
	w.Outdent()
	w.WriteType("non-indented message", Info)
}
