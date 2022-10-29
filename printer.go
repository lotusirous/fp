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
	// Write the rows.
	Write(w io.Writer) (n int, err error)
}

// NewRowRender creates the printer.
func NewRowRender() RowRenderer {
	return &pw{make([]row, 0)}
}

type row struct {
	header string
	value  string
}
type pw struct {
	rows []row
}

// AddRow adds a row to writer.
func (p *pw) AddRow(field, value string) {
	r := row{header: field, value: value}
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
		r := row{header: name, value: values[name]}
		p.rows = append(p.rows, r)
	}
}

// Render the rows to output
func (p *pw) Write(w io.Writer) (int, error) {
	for _, row := range p.rows {
		if row.header == "" {
			fmt.Fprintf(w, "%s\n", row.value)
		} else {
			fmt.Fprintf(w, "%s: %s\n", strings.ToUpper(row.header), row.value)
		}
	}
	return 0, nil
}
