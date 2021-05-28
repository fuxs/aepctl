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
	"io"
)

type TruncateWriter struct {
	w   io.Writer
	max int
	n   int
}

func NewTruncateWriter(w io.Writer, max int) *TruncateWriter {
	return &TruncateWriter{w: w, max: max}
}

func (t *TruncateWriter) Write(p []byte) (int, error) {
	var (
		b                    rune // current rune
		i, s, current, total int  // s rune start, i rune index, j bytes read
		err                  error
	)
	c := []rune(string(p))
	l := len(c)
	// are we still skipping?
	if t.n == t.max {
		for ; i < l; i++ {
			b = c[i]
			if b == '\n' || b == '\r' {
				s = i
				t.n = 0
				break
			}
		}
		if i == l {
			return len(p), nil
		}
		i++
	}
	// write lines
	for ; i < l; i++ {
		b = c[i]
		if b == '\n' || b == '\r' {
			t.n = 0
			continue
		}
		t.n++
		if t.n == t.max {
			if current, err = t.w.Write([]byte(string(c[s:i]))); err != nil {
				return total + current, err
			}
			total += current
			if _, err = t.w.Write([]byte("…")); err != nil {
				return total, err
			}
			for ; i < l; i++ {
				b = c[i]
				if b == '\n' || b == '\r' {
					t.n = 0
					break
				}
			}
			s = i
		}
	}
	// write the rest
	if s < l {
		if current, err = t.w.Write([]byte(string(c[s:l]))); err != nil {
			return total + current, err
		}
	}
	return len(p), nil
}
