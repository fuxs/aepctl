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
	"time"

	"gopkg.in/yaml.v3"
)

type TableDescriptor struct {
	Columns []*TableColumnDescriptor `json:"columns,omitempty" yaml:"columns,omitempty"`
	Path    []string                 `json:"path,omitempty" yaml:"path,omitempty"`
}

func NewTableDescriptor(def string) (*TableDescriptor, error) {
	result := &TableDescriptor{}
	if err := yaml.Unmarshal([]byte(def), &result); err != nil {
		return nil, err
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
	result := make([]string, len(t.Columns))
	for i, c := range t.Columns {
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
	out := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		if c.ID {
			out[i] = name
			continue
		}
		out[i] = c.Extract(q)
	}
	return w.Write(out...)
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
	switch t.Type {
	case "str":
		switch t.Format {
		case "":
			t.o = simpleString
			return
		}
	case "num":
		switch t.Format {
		case "utime":
			t.o = func(q *Query) string {
				v := q.Integer()
				if v == 0 {
					return "-"
				}
				return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
			}
			return
		}
	}
	t.o = simpleString
}
