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
	"runtime"
)

// Query supports queries on raw json objects
type Query struct {
	obj interface{}
	jp  jsonPath
}

// NewQuery creates an initialized query object
func NewQuery(obj interface{}) *Query {
	return &Query{obj: obj}
}

func NewQueryM(obj interface{}, jp jsonPath) *Query {
	return &Query{obj: obj, jp: jp}
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
	if len(q.jp) == 0 {
		runtime.Breakpoint()
	}
	return q.jp.Name()
}

func (q *Query) JSONPath() string {
	return q.jp.Path()
}

func (q *Query) JSONFullPath() string {
	return q.jp.String()
}

// Path queries nested objects, e.g. property a.b.c will be queried with Path("a","b","c")
func (q *Query) Path(path ...string) *Query {
	cur := q.obj
	jp := q.jp.Clone()
	for _, p := range path {
		if next, ok := cur.(map[string]interface{}); ok {
			cur = next[p]
			jp.Push(NewJSONPathAttribute(p))
		} else {
			return &Query{}
		}
	}
	return NewQueryM(cur, jp)
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
		for index, obj := range ar {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathIndex(index))
			response[index] = &Query{obj: obj, jp: jp}
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
		for index, obj := range ar {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathIndex(index))
			rf(&Query{obj: obj, jp: jp})
		}
	}
}

// RangeI executes the passed function on all children of the current object. It provides the index of the object.
func (q *Query) RangeI(rf func(int, *Query)) {
	if ar, ok := q.obj.([]interface{}); ok {
		for index, obj := range ar {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathIndex(index))
			rf(index, &Query{obj: obj, jp: jp})
		}
	}
}

// Concat calls the passed function on all children and concatenates the results separated by the passed separator.
func (q *Query) Concat(separator string, rf func(*Query) string) string {
	var buffer bytes.Buffer
	if ar, ok := q.obj.([]interface{}); ok {
		next := false
		for index, obj := range ar {
			if next {
				buffer.WriteString(separator)
			} else {
				next = true
			}
			jp := q.jp.Clone()
			jp.Push(NewJSONPathIndex(index))
			buffer.WriteString(rf(&Query{obj: obj, jp: jp}))
		}
	}
	return buffer.String()
}

// Get returns the children at the passed index
func (q *Query) Get(index int) *Query {
	var result *Query
	if ar, ok := q.obj.([]interface{}); ok {
		if index < len(ar) {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathIndex(index))
			result = &Query{obj: ar[index], jp: jp}
		}
	}
	return result
}

// RangeAttributes executes the passed function on all children of the current object
func (q *Query) RangeAttributes(rf func(string, *Query)) {
	if ar, ok := q.obj.(map[string]interface{}); ok {
		for k, v := range ar {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathAttribute(k))
			rf(k, &Query{obj: v, jp: jp})
		}
	}
}

// RangeAttributesE executes the passed function on all children of the current object
func (q *Query) RangeAttributesE(rf func(string, *Query) error) error {
	if ar, ok := q.obj.(map[string]interface{}); ok {
		for k, v := range ar {
			jp := q.jp.Clone()
			jp.Push(NewJSONPathAttribute(k))
			if err := rf(k, &Query{obj: v, jp: jp}); err != nil {
				return err
			}
		}
	}
	return nil
}
