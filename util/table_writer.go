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
	}
	return result
}

// NewCSVWriter creates an initialized RowWriter with commas as delimiter
func NewCSVWriter(out io.Writer) *RowWriter {
	result := &RowWriter{
		w: out,
		d: ",",
	}
	return result
}

// AutoFlush sets the limit for the automatic flush during writes
func (t *RowWriter) AutoFlush(l int) *RowWriter {
	t.l = l
	return t
}

// Write writes one row and terminates it with a newline
func (t *RowWriter) Write(v ...string) error {
	for i, w := range v {
		if i > 0 {
			// write delimiter
			if _, err := t.w.Write([]byte(t.d)); err != nil {
				return err
			}
		}
		// trim values, we have seen tabs in content
		if _, err := t.w.Write([]byte(strings.Trim(w, " \t"))); err != nil {
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

// Flush flushes the underlying stream
func (t *RowWriter) Flush() error {
	t.c = 0
	if t.f != nil {
		return t.f()
	}
	return nil
}
