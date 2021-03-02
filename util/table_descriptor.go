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
	"errors"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type TableDescriptor struct {
	Columns []*TableColumnDescriptor `json:"columns" yaml:"columns"`
	Wide    []*TableColumnDescriptor `json:"wide,omitempty" yaml:"wide,omitempty"`
	Path    []string                 `json:"path,omitempty" yaml:"path,omitempty"`
	Iter    string                   `json:"iterator,omitempty" yaml:"iterator,omitempty"`
	Filter  []string                 `json:"filter,omitempty" yaml:"filter,omitempty"`
}

func NewTableDescriptor(def string) (*TableDescriptor, error) {
	result := &TableDescriptor{}
	if err := yaml.Unmarshal([]byte(def), &result); err != nil {
		return nil, err
	}
	if result.Iter == "" {
		result.Iter = "array"
	}
	switch result.Iter {
	case "filter":
		if len(result.Filter) == 0 {
			return nil, errors.New("Iterator type filter requires filter attribute")
		}
	case "object", "array":
		if len(result.Filter) > 0 {
			return nil, fmt.Errorf("Iterator type %v does not support filter attribute", result.Iter)
		}
	}
	for i, c := range result.Columns {
		if c.Name == "" {
			return nil, fmt.Errorf("Name is empty in column %v", i)
		}
		if !c.ID {
			if c.Type == "" {
				return nil, fmt.Errorf("Type is empty in column %v", i)
			}
			if len(c.Path) == 0 {
				return nil, fmt.Errorf("Path is empty in column %v", i)
			}
		}
	}
	for i, c := range result.Wide {
		if c.Name == "" {
			return nil, fmt.Errorf("Name is empty in column %v", i)
		}
		if !c.ID {
			if c.Type == "" {
				return nil, fmt.Errorf("Type is empty in column %v", i)
			}
			if len(c.Path) == 0 {
				return nil, fmt.Errorf("Path is empty in column %v", i)
			}
		}
	}
	return result, nil
}

func ReadTableDescritorYAML(r io.Reader) (*TableDescriptor, error) {
	dec := yaml.NewDecoder(r)
	result := &TableDescriptor{}
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (t *TableDescriptor) Header(wide bool) []string {
	cols := t.Columns
	if wide && len(t.Wide) > 0 {
		cols = t.Wide
	}
	result := make([]string, len(cols))
	for i, c := range cols {
		result[i] = c.Name
	}
	return result
}

func (t *TableDescriptor) Preprocess(i JSONResponse) error {
	if len(t.Path) > 0 {
		if err := i.Path(t.Path...); err != nil {
			return err
		}
	}
	return i.Enter()
}

func (t *TableDescriptor) WriteRow(name string, q *Query, w *RowWriter, wide bool) error {
	cols := t.Columns
	if wide && len(t.Wide) > 0 {
		cols = t.Wide
	}
	out := make([]string, len(cols))
	for i, c := range cols {
		if c.ID {
			out[i] = name
			continue
		}
		out[i] = c.Extract(q)
	}
	return w.Write(out...)
}

func (t *TableDescriptor) Iterator(stream io.ReadCloser) (JSONResponse, error) {
	switch t.Iter {
	case "array":
		return NewJSONIterator(stream)
	case "filter":
		return NewJSONFilterIterator(t.Filter, stream)
	case "object":
		return NewJSONMapIterator(stream)
	default:
		return nil, fmt.Errorf("Unknown iterator %v", t.Iter)
	}
}

type TableColumnDescriptor struct {
	Name       string   `json:"name" yaml:"name"`
	Type       string   `json:"type" yaml:"type"`
	ID         bool     `json:"id,omitempty" yaml:"id,omitempty"`
	Path       []string `json:"path,omitempty" yaml:"path,omitempty"`
	Format     string   `json:"format,omitempty" yaml:"format,omitempty"`
	Parameters []string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	o          func(*Query) string
}

func simpleString(q *Query) string {
	return q.String()
}

func (t *TableColumnDescriptor) Extract(q *Query) string {
	r := q.Path(t.Path...)
	if t.o == nil {
		//t.completeType(r)
		t.assignFunc()
	}
	return t.o(r)
}

/*func (t *TableColumnDescriptor) completeType(q *Query) {
	if t.Type != "" {
		return
	}
	switch reflect.ValueOf(q.Interface()).Kind() {
	case reflect.String:
		t.Type = "string"
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Float32, reflect.Float64:
		t.Type = "number"
	}
}*/

func (t *TableColumnDescriptor) assignFunc() {
	t.o = simpleString
	switch t.Type {
	case "str":
		switch t.Format {
		case "":
			t.o = simpleString
			return
		}
	case "num":
		switch t.Format {
		case "":
			t.o = func(q *Query) string {
				return q.String()
			}
		case "utime":
			t.o = func(q *Query) string {
				v := q.Integer()
				if v == 0 {
					return "-"
				}
				return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
			}
		case "duration":
			t.o = func(q *Query) string {
				return time.Duration(q.Integer() * int(time.Millisecond)).String()
			}
		}
	}
}
