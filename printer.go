package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// RowRenderer render the row to a file descriptor.
type RowRenderer interface {
	// AddRow adds a row to writer.
	AddRow(field, value string)
	// AddRowMap appends the row by map input.
	AddRowMap(values map[string]string)
	// Render the rows.
	Render()
}

// NewRowRender creates the printer.
func NewRowRender(w io.Writer) RowRenderer {
	return &pw{w, make([][]string, 0)}
}

type pw struct {
	w    io.Writer
	rows [][]string
}

// AddRow adds a row to writer.
func (p *pw) AddRow(field, value string) {
	r := []string{field, value}
	p.rows = append(p.rows, r)
}

// AddRowMap appends the row from map. It sorts the map by default.
func (p *pw) AddRowMap(values map[string]string) {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		r := []string{name, values[name]}
		p.rows = append(p.rows, r)
	}

}

// Render the rows to output
func (p *pw) Render() {
	for _, l := range p.rows {
		fmt.Fprintf(p.w, "%s: %s\n", strings.ToUpper(l[0]), l[1])
	}
}
