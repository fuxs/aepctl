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
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// TableDescriptor contains all information to transform a JSON object to a table
type TableDescriptor struct {
	Columns []*TableColumnDescriptor `json:"columns,omitempty" yaml:"columns,omitempty"`
	//	Wide    []*TableColumnDescriptor `json:"wide,omitempty" yaml:"wide,omitempty"`
	Path   []string          `json:"path,omitempty" yaml:"path,omitempty"`
	Select []string          `json:"select,omitempty" yaml:"select,omitempty"`
	Iter   string            `json:"iterator,omitempty" yaml:"iterator,omitempty"`
	Filter []string          `json:"filter,omitempty" yaml:"filter,omitempty"`
	Vars   []*DescriptorVars `json:"vars,omitempty" yaml:"vars,omitempty"`
	Range  *DescriptorRange  `json:"range,omitempty" yaml:"range,omitempty"`
	thin   []*TableColumnDescriptor
	wide   []*TableColumnDescriptor
}

// NewTableDescriptor creates an initialzed TableDescriptor. It accpets a
// definition encoded in YAML format.
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
		return nil, errors.New("columns and range are defined")
	}
	var (
		l        int
		varTypes map[string]string
		cols     []*TableColumnDescriptor
	)
	if result.Range != nil {
		// range rows
		varTypes = make(map[string]string, len(result.Vars)+len(result.Range.Vars))
		for _, v := range result.Vars {
			varTypes[v.Name] = v.Type
		}
		for _, v := range result.Range.Vars {
			varTypes[v.Name] = v.Type
		}
		cols = result.Range.Columns
	} else {
		// simple column
		varTypes = make(map[string]string, len(result.Vars))
		for _, v := range result.Vars {
			varTypes[v.Name] = v.Type
		}
		cols = result.Columns
	}
	l = len(cols)
	w := make([]*TableColumnDescriptor, 0, l)
	t := make([]*TableColumnDescriptor, 0, l)
	for i, c := range cols {
		if c.Name == "" {
			return nil, fmt.Errorf("name is empty in column %v", i)
		}
		// determine column type
		if c.Meta == "" {
			// c.Type has highest priority
			if c.Type == "" {
				c.Type = "str"
				// if c.Var is set and has a type
				if c.Var != "" {
					if t := varTypes[c.Var]; t != "" {
						c.Type = t
					}
				}
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

// Header extracts the table header
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

// Preprocess goes down the path and enters the list or object
func (t *TableDescriptor) Preprocess(i JSONResponse) error {
	if len(t.Path) > 0 {
		if err := i.Path(t.Path...); err != nil {
			return err
		}
	}
	return i.Enter()
}

func processColumns(scope *Scope, cols []*TableColumnDescriptor, q *Query) []string {
	result := make([]string, len(cols))
	var value string
	for i, c := range cols {
		switch c.Meta {
		case "name":
			value = q.JSONName()
		case "path":
			value = q.JSONPath()
		default:
			value = c.Extract(scope, q)
		}
		result[i] = strings.Replace(value, "\t", " ", -1)
	}
	return result
}

// WriteRow writes one or more rows out
func (t *TableDescriptor) WriteRow(q *Query, w *RowWriter, wide bool) error {
	cols := t.thin
	if wide {
		cols = t.wide
	}
	if t.Range == nil {
		out := processColumns(rootScope, cols, q)
		return w.Write(out...)
	}
	r := t.Range
	s := NewScope(rootScope, t.Vars, q)
	return q.RangeAttributesE(func(name string, q *Query) error {
		ss := NewScope(s, r.Vars, q)
		out := processColumns(ss, cols, q)
		if r.Post != nil {
			for _, v := range r.Post.Vars {
				ss.Set(v.Name, v.Value)
			}
		}
		return w.Write(out...)
	})
}

// Iterator selects the configured iterator for the passed JSON stream
func (t *TableDescriptor) Iterator(stream io.ReadCloser) (JSONResponse, error) {
	switch t.Iter {
	case "array":
		return NewJSONIterator(stream), nil
	case "filter":
		return NewJSONFilterIterator(t.Filter, stream), nil
	case "object":
		return NewJSONMapIterator(stream), nil
	case "value":
		return NewJSONValueIterator(stream, t.Select), nil
	default:
		return nil, fmt.Errorf("unknown iterator %v", t.Iter)
	}
}

// StatusMapper maps status values to a pretty representation
var statusMapper = Mapper{
	"live":     "● Live",
	"approved": "● Approved",
	"draft":    "◯ Draft",
}

var stateMapper = Mapper{
	"enabled": "● Enabled",
}

// TableColumnDescriptor contains all information to extract a column value
type TableColumnDescriptor struct {
	Name       string   `json:"name" yaml:"name"`
	Long       string   `json:"long" yaml:"long"`
	Type       string   `json:"type" yaml:"type"`
	Meta       string   `json:"meta,omitempty" yaml:"meta,omitempty"`
	Path       []string `json:"path,omitempty" yaml:"path,omitempty"`
	Format     string   `json:"format,omitempty" yaml:"format,omitempty"`
	Parameters []string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Var        string   `json:"var,omitempty" yaml:"var,omitempty"`
	Mode       string   `json:"mode,omitempty" yaml:"mode,omitempty"`
	o          func(*Scope, *Query) string
}

// Extract retrieves the value from the JSON document and returns it as
// formatted string
func (t *TableColumnDescriptor) Extract(scope *Scope, q *Query) string {
	if t.o == nil {
		t.assignFunc()
	}
	if t.Var == "" {
		return t.o(scope, q.Path(t.Path...))
	}
	return t.o(scope, scope.Get(t.Var).Path(t.Path...))
}

func (t *TableColumnDescriptor) assignFunc() {
	t.o = func(_ *Scope, q *Query) string {
		return q.String()
	}
	switch t.Type {
	case "str":
		switch t.Format {
		case "localTime":
			if len(t.Parameters) > 0 {
				t.o = func(_ *Scope, q *Query) string {
					return LocalTimeStrCustom(q.String(), t.Parameters[0])
				}
			} else {
				t.o = func(_ *Scope, q *Query) string {
					return LocalTimeStr(q.String())
				}
			}
		case "status":
			t.o = func(_ *Scope, q *Query) string {
				return statusMapper.Lookup(q.String())
			}
		case "state":
			t.o = func(_ *Scope, q *Query) string {
				return stateMapper.Lookup(q.String())
			}
		}
	case "num":
		switch t.Format {
		case "utime":
			t.o = func(_ *Scope, q *Query) string {
				v := q.Integer()
				if v == 0 {
					return "-"
				}
				return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
			}
		case "duration":
			t.o = func(_ *Scope, q *Query) string {
				return time.Duration(q.Integer() * int(time.Millisecond)).String()
			}
		}
	case "list":
		switch t.Format {
		case "":
			t.o = func(_ *Scope, q *Query) string {
				return q.Concat(",", func(q *Query) string { return q.String() })
			}
		case "contains":
			t.o = func(_ *Scope, q *Query) string {
				return ContainsS(t.Parameters[0], q.Strings())
			}
		}
	}
}

// DescriptorVars represents a variable
type DescriptorVars struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
	Meta  string `json:"meta,omitempty" yaml:"meta,omitempty"`
	Cast  string `json:"cast,omitempty" yaml:"cast,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// DescriptorRange represents a range. Usually used to extract sub-values in
// multiple rows
type DescriptorRange struct {
	Vars    []*DescriptorVars        `json:"vars,omitempty" yaml:"vars,omitempty"`
	Columns []*TableColumnDescriptor `json:"columns,omitempty" yaml:"columns,omitempty"`
	Post    *RangePost               `json:"post,omitempty" yaml:"post,omitempty"`
}

// RangePost represents the post phase of a range, usually for variable update
type RangePost struct {
	Vars []*DescriptorVars `json:"vars,omitempty" yaml:"vars,omitempty"`
}

var rootScope = &Scope{vars: make(map[string]*Query)}

// Scope represents the current variable scope
type Scope struct {
	parent *Scope
	vars   map[string]*Query
}

// Get returns the value of a variable
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

// Set changes an existing variable in this or a parent scope or creates a new
// one
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

// NewScope creates an initialized Scope object
func NewScope(parent *Scope, vars []*DescriptorVars, q *Query) *Scope {
	result := make(map[string]*Query, len(vars))
	for _, v := range vars {
		switch v.Meta {
		case "name":
			result[v.Name] = NewQuery(q.JSONName())
			continue
		case "path":
			result[v.Name] = NewQuery(q.JSONPath())
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
