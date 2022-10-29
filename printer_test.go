package main

import (
	"bytes"
	"testing"
)

func TestRowPrinter(t *testing.T) {
	buf := new(bytes.Buffer)
	pw := &pw{make([]row, 0)}
	pw.AddRow("FILE", "bar")
	pw.Write(buf)

	got := buf.String()
	want := "FILE: bar\n"

	if got != want {
		t.Errorf("Got: %s - want: %s", buf, want)
	}
}
