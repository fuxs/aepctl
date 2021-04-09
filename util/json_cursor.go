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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type jsonState int

const (
	JSONS_UNDEFINED jsonState = iota
	JSONS_OPEN
	JSONS_DONE
	JSONS_O  // object attribute, expecting string or }
	JSONS_OV // object value, expecting {, [ VALUE or }
	JSONS_A  // array value, expecting {, [ VALUE or ]
	JSONS_AV // close array value, required for simple values
)

type jsonStateStack []jsonState

func (jss *jsonStateStack) Push(js jsonState) {
	*jss = append(*jss, js)
}

func (jss *jsonStateStack) Peek() jsonState {
	l := len(*jss)
	if l == 0 {
		return JSONS_UNDEFINED
	}
	return (*jss)[l-1]
}

func (jss *jsonStateStack) Pop() jsonState {
	l := len(*jss)
	if l == 0 {
		return JSONS_UNDEFINED
	}
	result := (*jss)[l-1]
	*jss = (*jss)[:l-1]
	return result
}

type jsonPathElement struct {
	array bool
	name  string
	index int
}

func NewJSONPathAttribute(name string) *jsonPathElement {
	return &jsonPathElement{name: name}
}

func NewJSONPathIndex(index int) *jsonPathElement {
	return &jsonPathElement{array: true, index: index, name: strconv.FormatInt(int64(index), 10)}
}

type jsonPath []*jsonPathElement

func (ps *jsonPath) Push(p ...*jsonPathElement) {
	*ps = append(*ps, p...)
}

func (ps *jsonPath) Peek() *jsonPathElement {
	l := len(*ps)
	if l == 0 {
		return nil
	}
	result := (*ps)[l-1]
	return result
}

func (ps *jsonPath) Pop() *jsonPathElement {
	l := len(*ps)
	if l == 0 {
		return nil
	}
	result := (*ps)[l-1]
	*ps = (*ps)[:l-1]
	return result
}

func (ps *jsonPath) Name() string {
	l := len(*ps)
	if l == 0 {
		return ""
	}
	last := (*ps)[l-1]
	if last.array {
		var buffer strings.Builder
		buffer.WriteRune('[')
		buffer.WriteString(last.name)
		buffer.WriteRune(']')
		return buffer.String()
	}
	return last.name
}

func (ps *jsonPath) Path() string {
	l := len(*ps)
	if l < 2 {
		return ""
	}
	var buffer strings.Builder
	next := false
	for _, e := range (*ps)[:l-1] {
		if next {
			if !e.array {
				buffer.WriteString(".")
			}
		} else {
			next = true
		}
		if e.array {
			buffer.WriteRune('[')
			buffer.WriteString(e.name)
			buffer.WriteRune(']')
		} else {
			buffer.WriteString(e.name)
		}
	}
	return buffer.String()
}

func (ps *jsonPath) String() string {
	var buffer strings.Builder
	next := false
	for _, e := range *ps {
		if next {
			if !e.array {
				buffer.WriteString(".")
			}
		} else {
			next = true
		}
		if e.array {
			buffer.WriteRune('[')
			buffer.WriteString(e.name)
			buffer.WriteRune(']')
		} else {
			buffer.WriteString(e.name)
		}
	}
	return buffer.String()
}

func (ps *jsonPath) Inc() {
	if e := ps.Peek(); e != nil {
		e.index++
		e.name = strconv.FormatInt(int64(e.index), 10)
	}
}

// Matches returns true,false
func (ps *jsonPath) Matches(path []string) bool {
	if len(path) > len(*ps) {
		return false
	}
	for i, pe := range path {
		if pe == "?" || pe == (*ps)[i].name {
			continue
		}
		return false
	}
	return true
}

func (ps *jsonPath) Clone() jsonPath {
	result := make(jsonPath, len(*ps))
	copy(result, *ps)
	return result
}

type JSONCursor struct {
	dec    *json.Decoder
	stream io.ReadCloser
	jss    jsonStateStack
	jp     jsonPath
}

func NewJSONCursor(stream io.ReadCloser) *JSONCursor {
	jss := make(jsonStateStack, 0, 16)
	jss.Push(JSONS_OPEN)
	jp := make(jsonPath, 0, 16)
	return &JSONCursor{dec: json.NewDecoder(stream), stream: stream, jss: jss, jp: jp}
}

func (j *JSONCursor) PathInfo() (string, string) {
	l := len(j.jp)
	if l == 0 {
		return "", ""
	}
	name := j.jp.Name()
	if l == 1 {
		return name, ""
	}
	return name, j.jp.Path()
}

func (j *JSONCursor) More() bool {
	return j.dec.More()
}

func (j *JSONCursor) Offset() int64 {
	return j.dec.InputOffset()
}

func (j *JSONCursor) MoreTokens() bool {
	return j.jss.Peek() != JSONS_DONE
}

func (j *JSONCursor) NextValueF(filter []string) (*Query, error) {
	q, err := j.NextValue()
	for err == nil {
		if q.jp.Matches(filter) {
			return q, nil
		}
		q, err = j.NextValue()
	}
	return nil, err
}

func (j *JSONCursor) NextValue() (*Query, error) {
	state := j.jss.Peek()
	for state != JSONS_DONE {
		path := j.jp.Clone()
		t, err := j.Token()
		if err != nil {
			return nil, err
		}
		if state == JSONS_OV || state == JSONS_A || state == JSONS_AV {
			if _, ok := t.(json.Delim); !ok {
				return NewQueryM(t, path), nil
			}
		}
		state = j.jss.Peek()
	}
	return nil, io.EOF
}

func (j *JSONCursor) Token() (json.Token, error) {
	state := j.jss.Peek()
	if state == JSONS_DONE {
		return nil, io.EOF
	}
	if state == JSONS_AV {
		j.jp.Inc()
		j.jss.Pop()
		state = j.jss.Peek()
	}
	t, err := j.dec.Token()
	if err != nil {
		return nil, err
	}
	switch state {
	case JSONS_OPEN:
		d, ok := t.(json.Delim)
		if !ok || !(d == '{' || d == '[') {
			return nil, fmt.Errorf("expecting [ or { at position %v", j.dec.InputOffset())
		}
		j.jss.Pop()
		j.jss.Push(JSONS_DONE)
		if d == '{' {
			j.jss.Push(JSONS_O)
		} else {
			j.jss.Push(JSONS_A)
		}
	case JSONS_O:
		str, ok := t.(string)
		if ok {
			j.jp.Push(NewJSONPathAttribute(str))
			j.jss.Push(JSONS_OV)
			return t, nil
		}
		d, ok := t.(json.Delim)
		if !ok || !(d == '}') {
			return nil, fmt.Errorf("expecting } position %v", j.dec.InputOffset())
		}
		j.jss.Pop()
		if j.jss.Peek() == JSONS_A {
			j.jp.Inc()
		} else {
			j.jp.Pop()
		}
	case JSONS_OV:
		j.jss.Pop()
		d, ok := t.(json.Delim)
		if ok {
			switch d {
			case '{':
				j.jss.Push(JSONS_O)
			case '[':
				j.jss.Push(JSONS_A)
				j.jp.Push(NewJSONPathIndex(0))
			case '}':
				if j.jss.Peek() == JSONS_A {
					// increment array index
					j.jp.Inc()
				} else {
					// must be object, pop attribute name
					j.jp.Pop()
				}
			default:
				return nil, fmt.Errorf("expecting [,{ or } at position %v", j.dec.InputOffset())
			}
		} else {
			// it's a value
			j.jp.Pop()
		}
	case JSONS_A:
		d, ok := t.(json.Delim)
		if ok {
			switch d {
			case '{':
				j.jss.Push(JSONS_O)
			case '[':
				j.jss.Push(JSONS_A)
				j.jp.Push(NewJSONPathIndex(0))
			case ']':
				j.jss.Pop()
				j.jp.Pop()
				if j.jss.Peek() == JSONS_A {
					j.jp.Inc()
				} else {
					j.jp.Pop()
				}
			default:
				return nil, fmt.Errorf("expecting [,{ or ] at position %v", j.dec.InputOffset())
			}
		} else {
			j.jss.Push(JSONS_AV)
		}
	default:
		return nil, errors.New("state error")
	}
	return t, nil
}

// Skip skips the next element like for example a string, object or array
func (j *JSONCursor) Skip() error {
	state := j.jss.Peek()
	if state == JSONS_DONE {
		return io.EOF
	}
	t, err := j.Token()
	if err != nil {
		return err
	}
	d, ok := t.(json.Delim)
	if !ok || d == '}' || d == ']' {
		// everything is fine
		return nil
	}
	counter := 1
	if d == '{' {
		t, err = j.Token()
		for err == nil {
			d, ok = t.(json.Delim)
			if ok {
				switch d {
				case '{':
					counter++
				case '}':
					counter--
					if counter == 0 {
						return nil
					}
				}
			}
			t, err = j.Token()
		}
		return err
	}
	// must be '['
	t, err = j.Token()
	for err == nil {
		d, ok = t.(json.Delim)
		if ok {
			switch d {
			case '[':
				counter++
			case ']':
				counter--
				if counter == 0 {
					return nil
				}
			}
		}
		t, err = j.Token()
	}
	return err
}

func (j *JSONCursor) Decode(v interface{}) error {
	state := j.jss.Peek()
	if state == JSONS_DONE {
		return io.EOF
	}
	if state != JSONS_OV && state != JSONS_A && state != JSONS_OPEN {
		return errors.New("state error")
	}
	if err := j.dec.Decode(v); err != nil {
		return err
	}
	if state != JSONS_A {
		j.jss.Pop()
	}
	if state == JSONS_O {
		j.jp.Pop()
	}
	return nil
}

// Close closes the underlying ReaderCloser stream
func (j *JSONCursor) Close() error {
	return j.stream.Close()
}

// PrintRaw copies the raw data to standard out
func (j *JSONCursor) PrintRaw() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	_, err := io.Copy(bout, j.stream)
	return err
}

// PrintPretty prints the raw data with indention to standard out
func (j *JSONCursor) PrintPretty() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	return JSONPrintPretty(j.dec, bout)
}
