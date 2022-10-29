package main

import (
	"bytes"
	"testing"
)

func TestRowPrinter(t *testing.T) {
	buf := new(bytes.Buffer)
	pw := &pw{buf, make([][]string, 0)}
	pw.AddRow("FILE", "bar")
	pw.RenderTo(buf)

	got := buf.String()
	want := "FILE: bar\n"

	if got != want {
		t.Errorf("Got: %s - want: %s", buf, want)
	}
}
