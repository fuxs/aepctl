/*
Package util util consists of general utility functions and structures.

Copyright 2021 Michael Bungenstock

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package util

import (
	"fmt"
	"io"

	"text/tabwriter"
)

// Table contains any object represented as table
type Table struct {
	Columns    []*Column
	Height     int
	Width      int
	HideHeader bool
	columns    map[string]*Column
}

// Column represents one column in a table
type Column struct {
	Number int
	Name   string
	Cells  []interface{}
	Calc   func(interface{}) interface{}
}

// NewTable creates a new initialized table with column names and initial capacity
func NewTable(names []string, capacity int) *Table {
	cols := len(names)
	result := &Table{
		Columns: make([]*Column, cols),
		Height:  0,
		Width:   cols,
		columns: make(map[string]*Column, cols),
	}
	for i := 0; i < cols; i++ {
		c := &Column{
			Name:   names[i],
			Number: i,
			Cells:  make([]interface{}, 0, capacity),
		}
		result.Columns[i] = c
		result.columns[names[i]] = c
	}
	return result
}

// Append appends a new row to the end of the table
func (t *Table) Append(row map[string]interface{}) {
	for k, column := range t.columns {
		column.Cells = append(column.Cells, row[k])
	}
	t.Height++
}

// Column returns the column with the passed name
func (t *Table) Column(name string) *Column {
	return t.columns[name]
}

// Get returns the cell value for the passed coordinations
func (t *Table) Get(row, col int) interface{} {
	c := t.Columns[col]
	if c.Calc != nil {
		return c.Calc(c.Cells[row])
	}
	return c.Cells[row]
}

// Array converts the table to a two-dimensional array (row x column)
func (t *Table) Array() [][]interface{} {
	result := make([][]interface{}, t.Height+1)
	row := make([]interface{}, t.Width)
	for c := 0; c < t.Width; c++ {
		row[c] = t.Columns[c].Name
	}
	result[0] = row
	if t.Height > 0 {
		rows := result[1:]
		for r := 0; r < t.Height; r++ {
			row := make([]interface{}, t.Width)
			for c := 0; c < t.Width; c++ {
				row[c] = t.Get(r, c)
			}
			rows[r] = row
		}
	}
	return result
}

// Print prints the table to the passed io.Writer
func (t *Table) Print(out io.Writer) {
	w := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)
	if !t.HideHeader {
		for c := 0; c < t.Width; c++ {
			// ignore errors
			_, _ = w.Write([]byte(t.Columns[c].Name))
			_, _ = w.Write([]byte("\t"))
		}
		_, _ = w.Write([]byte("\n"))
	}
	for r := 0; r < t.Height; r++ {
		for c := 0; c < t.Width; c++ {
			// ignore errors
			_, _ = fmt.Fprint(w, t.Get(r, c))
			_, _ = w.Write([]byte("\t"))
		}
		_, _ = w.Write([]byte("\n"))
	}
	w.Flush()
}
