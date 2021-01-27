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

import "bytes"

// Query supports queries on raw json objects
type Query struct {
	obj interface{}
}

// NewQuery creates an initialized query object
func NewQuery(obj interface{}) *Query {
	return &Query{obj: obj}
}

// Path queries nested objects, e.g. property a.b.c will be queried with Path("a","b","c")
func (q *Query) Path(path ...string) *Query {
	cur := q.obj
	for _, p := range path {
		if next, ok := cur.(map[string]interface{}); ok {
			cur = next[p]
		} else {
			return &Query{obj: nil}
		}
	}
	return &Query{obj: cur}
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
func (q *Query) Range(rf func(*Query)) {
	if ar, ok := q.obj.([]interface{}); ok {
		for _, obj := range ar {
			rf(&Query{obj: obj})
		}
	}
}

// RangeI executes the passed function on all children of the current object. It provides the index of the object.
func (q *Query) RangeI(rf func(int, *Query)) {
	if ar, ok := q.obj.([]interface{}); ok {
		for i, obj := range ar {
			rf(i, &Query{obj: obj})
		}
	}
}

// Concat calls the passed function on all children and concatenates the results separated by the passed separator.
func (q *Query) Concat(separator string, rf func(*Query) string) string {
	var buffer bytes.Buffer
	if ar, ok := q.obj.([]interface{}); ok {
		next := false
		for _, obj := range ar {
			if next {
				buffer.WriteString(separator)
			} else {
				next = true
			}
			buffer.WriteString(rf(&Query{obj: obj}))

		}
	}
	return buffer.String()
}

// Get returns the children at the passed index
func (q *Query) Get(index int) *Query {
	var result *Query
	if ar, ok := q.obj.([]interface{}); ok {
		if index < len(ar) {
			result = &Query{obj: ar[index]}
		}
	}
	return result
}

// RangeAttributes executes the passed function on all children of the current object
func (q *Query) RangeAttributes(rf func(string, *Query)) {
	if ar, ok := q.obj.(map[string]interface{}); ok {
		for k, v := range ar {
			rf(k, &Query{obj: v})
		}
	}
}
