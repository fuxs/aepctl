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
	Columns []*TableColumnDescriptor `json:"columns,omitempty" yaml:"columns,omitempty"`
	//	Wide    []*TableColumnDescriptor `json:"wide,omitempty" yaml:"wide,omitempty"`
	Path   []string          `json:"path,omitempty" yaml:"path,omitempty"`
	Iter   string            `json:"iterator,omitempty" yaml:"iterator,omitempty"`
	Filter []string          `json:"filter,omitempty" yaml:"filter,omitempty"`
	Vars   []*DescriptorVars `json:"vars,omitempty" yaml:"vars,omitempty"`
	Range  *DescriptorRange  `json:"range,omitempty" yaml:"range,omitempty"`
	thin   []*TableColumnDescriptor
	wide   []*TableColumnDescriptor
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
	if len(result.Columns) > 0 && result.Range != nil {
		return nil, errors.New("Columns and range are defined")
	}
	var (
		l     int
		types map[string]string
		cols  []*TableColumnDescriptor
	)
	if result.Range != nil {
		types = make(map[string]string, len(result.Vars)+len(result.Range.Vars))
		for _, v := range result.Vars {
			types[v.Name] = v.Type
		}
		for _, v := range result.Range.Vars {
			types[v.Name] = v.Type
		}
		cols = result.Range.Columns
	} else {
		types = make(map[string]string, len(result.Vars))
		for _, v := range result.Vars {
			types[v.Name] = v.Type
		}
		cols = result.Columns
	}
	l = len(cols)
	w := make([]*TableColumnDescriptor, 0, l)
	t := make([]*TableColumnDescriptor, 0, l)
	for i, c := range cols {
		if c.Name == "" {
			return nil, fmt.Errorf("Name is empty in column %v", i)
		}
		if !c.ID {
			if c.Var != "" {
				c.Type = types[c.Var]
			} else if c.Type == "" {
				c.Type = "str"
			}
		}
		switch c.Mode {
		case "thin":
			t = append(t, c)
		case "wide":
			w = append(w, c)
		default:
			t = append(t, c)
			w = append(w, c)
		}
	}
	result.thin = t
	result.wide = w
	return result, nil
}

func ReadTableDescriptorYAML(r io.Reader) (*TableDescriptor, error) {
	dec := yaml.NewDecoder(r)
	result := &TableDescriptor{}
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (t *TableDescriptor) Header(wide bool) []string {
	cols := t.thin
	if wide {
		cols = t.wide
	}
	result := make([]string, len(cols))
	for i, c := range cols {
		if wide && c.Long != "" {
			result[i] = c.Long
		} else {
			result[i] = c.Name
		}

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

func processColumns(scope *Scope, cols []*TableColumnDescriptor, name string, q *Query) []string {
	result := make([]string, len(cols))
	for i, c := range cols {
		if c.ID {
			result[i] = name
			continue
		}
		result[i] = c.Extract(scope, q)
	}
	return result
}

func (t *TableDescriptor) WriteRow(name string, q *Query, w *RowWriter, wide bool) error {
	cols := t.thin
	if wide {
		cols = t.wide
	}
	if t.Range == nil {
		out := processColumns(rootScope, cols, name, q)
		return w.Write(out...)
	}
	r := t.Range
	s := NewScope(rootScope, t.Vars, name, q)
	q.RangeAttributesE(func(name string, q *Query) error {
		ss := NewScope(s, r.Vars, name, q)
		out := processColumns(ss, cols, name, q)
		if r.Post != nil {
			for _, v := range r.Post.Vars {
				ss.Set(v.Name, v.Value)
			}
		}
		return w.Write(out...)
	})
	return nil
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
	Long       string   `json:"long" yaml:"long"`
	Type       string   `json:"type" yaml:"type"`
	ID         bool     `json:"id,omitempty" yaml:"id,omitempty"`
	Path       []string `json:"path,omitempty" yaml:"path,omitempty"`
	Format     string   `json:"format,omitempty" yaml:"format,omitempty"`
	Parameters []string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Var        string   `json:"var,omitempty" yaml:"var,omitempty"`
	Mode       string   `json:"mode,omitempty" yaml:"mode,omitempty"`
	o          func(*Scope, *Query) string
}

func (t *TableColumnDescriptor) Extract(scope *Scope, q *Query) string {
	if t.o == nil {
		t.assignFunc()
	}
	return t.o(scope, q.Path(t.Path...))
}

func (t *TableColumnDescriptor) assignFunc() {
	vn := t.Var
	if vn == "" {
		t.o = func(_ *Scope, q *Query) string {
			return q.String()
		}
	} else {
		t.o = func(scope *Scope, _ *Query) string {
			return scope.Get(vn).String()
		}

	}
	switch t.Type {
	case "str":
		switch t.Format {
		case "localTime":
			if vn == "" {
				t.o = func(_ *Scope, q *Query) string {
					return LocalTimeStr(q.String())
				}
			} else {
				t.o = func(scope *Scope, _ *Query) string {
					q := scope.Get(vn)
					return LocalTimeStr(q.String())
				}
			}
		}
	case "num":
		switch t.Format {
		case "utime":
			if vn == "" {
				t.o = func(_ *Scope, q *Query) string {
					v := q.Integer()
					if v == 0 {
						return "-"
					}
					return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
				}
			} else {
				t.o = func(scope *Scope, _ *Query) string {
					v := scope.Get(vn).Integer()
					if v == 0 {
						return "-"
					}
					return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
				}
			}
		case "duration":
			if vn == "" {
				t.o = func(_ *Scope, q *Query) string {
					return time.Duration(q.Integer() * int(time.Millisecond)).String()
				}
			} else {
				t.o = func(scope *Scope, _ *Query) string {
					q := scope.Get(vn)
					return time.Duration(q.Integer() * int(time.Millisecond)).String()
				}
			}
		}
	case "list":
		switch t.Format {
		case "":
			if vn == "" {
				t.o = func(_ *Scope, q *Query) string {
					return q.Concat(",", func(q *Query) string { return q.String() })
				}
			} else {
				t.o = func(scope *Scope, _ *Query) string {
					q := scope.Get(vn)
					return q.Concat(",", func(q *Query) string { return q.String() })
				}
			}
		case "contains":
			if vn == "" {
				t.o = func(_ *Scope, q *Query) string {
					return ContainsS(t.Parameters[0], q.Strings())
				}
			} else {
				t.o = func(scope *Scope, _ *Query) string {
					q := scope.Get(vn)
					return ContainsS(t.Parameters[0], q.Strings())
				}
			}
		}
	}
}

type DescriptorVars struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
	ID    bool   `json:"id,omitempty" yaml:"id,omitempty"`
	Cast  string `json:"cast,omitempty" yaml:"cast,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type DescriptorRange struct {
	Vars    []*DescriptorVars        `json:"vars,omitempty" yaml:"vars,omitempty"`
	Columns []*TableColumnDescriptor `json:"columns,omitempty" yaml:"columns,omitempty"`
	Post    *RangePost               `json:"post,omitempty" yaml:"post,omitempty"`
}

type RangePost struct {
	Vars []*DescriptorVars `json:"vars,omitempty" yaml:"vars,omitempty"`
}

var rootScope = &Scope{vars: make(map[string]*Query)}

type Scope struct {
	parent *Scope
	vars   map[string]*Query
}

func (s *Scope) Get(name string) *Query {
	if v, found := s.vars[name]; found {
		return v
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return &Query{}
}

func (s *Scope) set(name, value string) bool {
	if _, found := s.vars[name]; found {
		s.vars[name] = NewQuery(value)
		return true
	}
	if s.parent != nil {
		return s.parent.set(name, value)
	}
	return false
}

func (s *Scope) Set(name, value string) bool {
	_, found := s.vars[name]
	if found {
		s.vars[name] = NewQuery(value)
		return true
	}
	if s.parent != nil && s.parent.set(name, value) {
		return true
	}
	s.vars[name] = NewQuery(value)
	return true
}

func NewScope(parent *Scope, vars []*DescriptorVars, name string, q *Query) *Scope {
	result := make(map[string]*Query, len(vars))
	for _, v := range vars {
		if v.ID {
			result[v.Name] = NewQuery(name)
			continue
		}
		switch v.Cast {
		case "strings":
			result[v.Name] = NewQuery(q.Strings())
			continue
		}
	}
	return &Scope{
		parent: parent,
		vars:   result,
	}
}
