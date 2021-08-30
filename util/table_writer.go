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
	"strings"
	"text/tabwriter"
)

// RowWriter writes rows to a stream
type RowWriter struct {
	w io.Writer    // writer
	f func() error // flush
	d string       // delimiter
	e string       // escaped delimiter
	c int          // counter
	l int          // limit
}

// NewTableWriter creates an initialized RowWriter with tabs as delimiter
func NewTableWriter(out io.Writer) *RowWriter {
	tw := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)
	result := &RowWriter{
		w: tw,
		f: tw.Flush,
		d: "\t",
		e: "    ",
	}
	return result
}

// NewCSVWriter creates an initialized RowWriter with commas as delimiter
func NewCSVWriter(out io.Writer) *RowWriter {
	result := &RowWriter{
		w: out,
		d: ",",
		e: ";",
	}
	return result
}

// AutoFlush sets the limit for the automatic flush during writes
func (t *RowWriter) AutoFlush(l int) *RowWriter {
	t.l = l
	return t
}

// Write writes one row and terminates it with a newline
func (t *RowWriter) WriteSingle(v ...string) error {
	for i, w := range v {
		if i > 0 {
			// write delimiter
			if _, err := t.w.Write([]byte(t.d)); err != nil {
				return err
			}
		}
		// trim values, we have seen tabs in content
		if _, err := t.w.Write([]byte(strings.Trim(w, "\t"))); err != nil {
			return err
		}
	}
	fmt.Fprintln(t.w)
	t.c++
	if t.l > 0 && t.c > t.l {
		if err := t.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Write writes one row with mutliple lines and terminates it with a newline. v
// is a slice of columns separated by the delimiter, e.g. a tab.
func (t *RowWriter) Write(v ...string) error {
	// l is number of columns
	l := len(v)
	values := make([]string, l)
	lengths := make([]int, l)
	// replace delimiters with escaped variant, e.g. tabs with spaces
	for i, w := range v {
		values[i] = strings.ReplaceAll(w, t.d, t.e)
		lengths[i] = len(values[i])
	}
	// stores the progress for each column
	offsets := make([]int, l)
	done := 0
	// iterate over all columns until each column has been printed row by row
	for done < l {
		// i is current column
		for i, w := range values {
			// write the column delimiter between columns
			if i > 0 {
				// write delimiter
				if _, err := t.w.Write([]byte(t.d)); err != nil {
					return err
				}
			}
			// check the length, skip empty
			length := lengths[i]
			if length == 0 {
				done++
				continue
			}
			offset := offsets[i]
			if offset < length {
				var str string
				current := w[offset:]
				// search for next newline
				index := strings.IndexByte(current, '\n')
				if index >= 0 {
					// skip newline characters
					skip := 1
					if index > 0 {
						if w[index-1] == '\r' {
							index--
							skip++
						}
					}
					str = current[:index]
					nextOffset := offset + index + skip
					if nextOffset == length {
						done++
					}
					offsets[i] = nextOffset
				} else {
					// no newline
					str = current
					offsets[i] = length
					done++
				}
				if _, err := t.w.Write([]byte(str)); err != nil {
					return err
				}
			}
		}
		// output new line
		fmt.Fprintln(t.w)
		// automatic flush
		t.c++
		if t.l > 0 && t.c > t.l {
			if err := t.Flush(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Flush flushes the underlying stream
func (t *RowWriter) Flush() error {
	t.c = 0
	if t.f != nil {
		return t.f()
	}
	return nil
}
