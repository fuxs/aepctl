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
	"bytes"
	"encoding/json"
	"io"
	"strconv"
)

// Query supports queries on raw json objects
type Query struct {
	obj  interface{}
	name string
	path string
}

// NewQuery creates an initialized query object
func NewQuery(obj interface{}) *Query {
	return &Query{obj: obj}
}

func NewQueryM(obj interface{}, name, path string) *Query {
	return &Query{obj: obj, name: name, path: path}
}

func NewQueryStream(stream io.Reader) (*Query, error) {
	dec := json.NewDecoder(stream)
	var obj interface{}
	if err := dec.Decode(&obj); err != nil {
		return nil, err
	}
	return NewQuery(obj), nil
}

// UnmarshallQuery unmarshal JSON data and returns a Query object
func UnmarshallQuery(data []byte) (*Query, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return NewQuery(obj), nil
}

func (q *Query) JSONName() string {
	return q.name
}

func (q *Query) JSONPath() string {
	return q.path
}

func (q *Query) JSONFullPath() string {
	return Concat(".", q.path, q.name)
}

// Path queries nested objects, e.g. property a.b.c will be queried with Path("a","b","c")
func (q *Query) Path(path ...string) *Query {
	cur := q.obj
	for _, p := range path {
		if next, ok := cur.(map[string]interface{}); ok {
			cur = next[p]
		} else {
			return &Query{}
		}
	}
	l := len(path)
	var name string
	if l == 0 {
		name = q.name
	} else {
		name = path[l-1]
	}
	var p string
	if l > 1 {
		p = Concat(".", path[:l-2]...)
	}
	return &Query{obj: cur, name: name, path: Concat(".", p)}
}

// Interface returns the current object
func (q *Query) Interface() interface{} {
	return q.obj
}

// Value returns the object of the referenced path
func (q *Query) Value(path ...string) interface{} {
	return q.Path(path...).obj
}

// Int returns the integer of the referenced path
func (q *Query) Int(path ...string) int {
	return GetInt(q.Path(path...).obj)
}

// Integer returns the current object as integer
func (q *Query) Integer() int {
	return GetInt(q.obj)
}

// Str returns the string of the referenced path
func (q *Query) Str(path ...string) string {
	return q.Path(path...).String()
}

// String returns the current object as string
func (q *Query) String() string {
	return GetString(q.obj)
}

// Strings returns the current object as array of strings
func (q *Query) Strings() []string {
	if ar, ok := q.obj.([]string); ok {
		return ar
	}
	if ar, ok := q.obj.([]interface{}); ok {
		result := make([]string, len(ar))
		for i, obj := range ar {
			result[i] = GetString(obj)
		}
		return result
	}
	return []string{}
}

// Len returns the length of the referenced path
func (q *Query) Len(path ...string) int {
	return q.Path(path...).Length()
}

// Length returns the length of the current object
func (q *Query) Length() int {
	if ar, ok := q.obj.([]interface{}); ok {
		return len(ar)
	}
	return 0
}

// Range executes the passed function on all children of the current object
func (q *Query) QueryArray() []*Query {
	if ar, ok := q.obj.([]interface{}); ok {
		response := make([]*Query, len(ar))
		path := q.JSONFullPath()
		for index, obj := range ar {
			name := strconv.FormatInt(int64(index), 10)
			response[index] = &Query{obj: obj, name: name, path: path}
		}
		return response
	}
	return []*Query{}
}

func (q *Query) Array() []interface{} {
	if ar, ok := q.obj.([]interface{}); ok {
		return ar
	}
	return []interface{}{}
}

// Range executes the passed function on all children of the current object
func (q *Query) Range(rf func(*Query)) {
	if ar, ok := q.obj.([]interface{}); ok {
		path := q.JSONFullPath()
		for index, obj := range ar {
			name := strconv.FormatInt(int64(index), 10)
			rf(&Query{obj: obj, name: name, path: path})
		}
	}
}

// RangeI executes the passed function on all children of the current object. It provides the index of the object.
func (q *Query) RangeI(rf func(int, *Query)) {
	if ar, ok := q.obj.([]interface{}); ok {
		path := q.JSONFullPath()
		for index, obj := range ar {
			name := strconv.FormatInt(int64(index), 10)
			rf(index, &Query{obj: obj, name: name, path: path})
		}
	}
}

// Concat calls the passed function on all children and concatenates the results separated by the passed separator.
func (q *Query) Concat(separator string, rf func(*Query) string) string {
	var buffer bytes.Buffer
	if ar, ok := q.obj.([]interface{}); ok {
		next := false
		path := q.JSONFullPath()
		for index, obj := range ar {
			if next {
				buffer.WriteString(separator)
			} else {
				next = true
			}
			name := strconv.FormatInt(int64(index), 10)
			buffer.WriteString(rf(&Query{obj: obj, name: name, path: path}))
		}
	}
	return buffer.String()
}

// Get returns the children at the passed index
func (q *Query) Get(index int) *Query {
	var result *Query
	if ar, ok := q.obj.([]interface{}); ok {
		if index < len(ar) {
			name := strconv.FormatInt(int64(index), 10)
			result = &Query{obj: ar[index], name: name, path: q.JSONFullPath()}
		}
	}
	return result
}

// RangeAttributes executes the passed function on all children of the current object
func (q *Query) RangeAttributes(rf func(string, *Query)) {
	if ar, ok := q.obj.(map[string]interface{}); ok {
		path := q.JSONFullPath()
		for k, v := range ar {
			rf(k, &Query{obj: v, name: k, path: path})
		}
	}
}

// RangeAttributesE executes the passed function on all children of the current object
func (q *Query) RangeAttributesE(rf func(string, *Query) error) error {
	if ar, ok := q.obj.(map[string]interface{}); ok {
		path := q.JSONFullPath()
		for k, v := range ar {
			if err := rf(k, &Query{obj: v, name: k, path: path}); err != nil {
				return err
			}
		}
	}
	return nil
}
