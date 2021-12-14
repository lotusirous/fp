package main

import (
	"bytes"
	"testing"
)

func TestRowPrinter(t *testing.T) {
	buf := new(bytes.Buffer)
	pw := &pw{buf, make([][]string, 0)}
	pw.AddRow("foo", "bar")
	pw.Render()

	got := buf.String()
	want := "FOO: bar\n"

	if got != want {
		t.Errorf("Got: %s - want: %s", buf, want)
	}

}
