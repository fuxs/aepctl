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
	"encoding/json"
	"fmt"
	"io"
)

// JSONResponse is the interface for streaming JSON objects
type JSONResponse interface {
	Close() error
	//Delim() (json.Delim, error)
	Enter() error
	EnterArray() error
	EnterObject() error
	Leave() error
	More() bool
	Next() (*Query, error)
	//Obj() (map[string]interface{}, error)
	Offset() int64
	Path(...string) error
	PrintRaw() error
	PrintPretty() error
	Query() (*Query, error)
	Range(func(*Query) error) error
	Skip() error
	Token() (json.Token, error)
}

// JSONIterator implements the JSONResponse interface
type JSONIterator struct {
	c *JSONCursor
	//dec    *json.Decoder
	//stream io.ReadCloser
}

// NewJSONIterator creates an initialized JSONIterator object
func NewJSONIterator(stream io.ReadCloser) *JSONIterator {
	return &JSONIterator{c: NewJSONCursor(stream)}
}

// More checks if there is another element in the current object or array
func (j *JSONIterator) More() bool {
	return j.c.More()
}

func (j *JSONIterator) Token() (json.Token, error) {
	result, err := j.c.Token()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (j *JSONIterator) Offset() int64 {
	return j.c.Offset()
}

// Skip skips the next element like for example a string, object or array
func (j *JSONIterator) Skip() error {
	return j.c.Skip()
}

// Enter moves forward to the first element
func (j *JSONIterator) Enter() error {
	d, err := j.Delim()
	if err != nil {
		return err
	}
	if d != '[' && d != '{' {
		return fmt.Errorf("expecting [ or { but found %v at offset %v", d, j.c.Offset())
	}
	return nil
}

// Enter moves forward to the first element
func (j *JSONIterator) Leave() error {
	d, err := j.Delim()
	if err != nil {
		return err
	}
	if d != ']' && d != '}' {
		return fmt.Errorf("expecting ] or } but found %v at offset %v", d, j.c.Offset())
	}
	return nil
}

// EnterArray moves forward to the first element in the array
func (j *JSONIterator) EnterArray() error {
	d, err := j.Delim()
	if err != nil {
		return err
	}
	if d != '[' {
		return fmt.Errorf("expecting [ found %v at offset %v", d, j.c.Offset())
	}
	return nil
}

// EnterObject moves forward to the first element in the object
func (j *JSONIterator) EnterObject() error {
	d, err := j.Delim()
	if err != nil {
		return err
	}
	if d != '{' {
		return fmt.Errorf("expected { at %v but found %v", j.c.Offset(), d)
	}
	return nil
}

// Path skips all element along the path
func (j *JSONIterator) Path(path ...string) error {
	for _, p := range path {
		if err := j.EnterObject(); err != nil {
			return err
		}
		found := false
		for j.More() {
			id, err := j.String()
			if err != nil {
				return err
			}
			if id == p || p == "?" {
				found = true
				break
			}
			if err = j.Skip(); err != nil {
				return err
			}
		}
		if !found {
			return fmt.Errorf("could not find id %v but reached end at offset %v", p, j.c.Offset())
		}
	}
	return nil
}

// Delim returns the delimiter like [,],{ or} at the current position
func (j *JSONIterator) Delim() (json.Delim, error) {
	t, err := j.c.Token()
	if err != nil {
		return 0, err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return 0, fmt.Errorf("parse error, expected delimiter at offset %v", j.c.Offset())
	}
	return delim, nil
}

// Obj returns the object at the current position
/*func (j *JSONIterator) Obj() (map[string]interface{}, error) {
	if !j.More() {
		return nil, fmt.Errorf("no more object available in current array or object at offset %v", j.c.Offset())
	}
	var obj map[string]interface{}
	if err := j.c.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}*/

// Interface returns any element at the current position
func (j *JSONIterator) Interface() (interface{}, error) {
	if !j.More() {
		return nil, fmt.Errorf("no more element available in current array or object at offset %v", j.c.Offset())
	}
	var obj interface{}
	if err := j.c.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (j *JSONIterator) Query() (*Query, error) {
	jp := j.c.jp.Clone()
	obj, err := j.Interface()
	if err != nil {
		return nil, err
	}
	return NewQueryM(obj, jp), nil
}

// Bool returns the boolean at the current position
func (j *JSONIterator) Bool() (bool, error) {
	if !j.More() {
		return false, fmt.Errorf("no more bool available in current array or object at offset %v", j.c.Offset())
	}
	t, err := j.c.Token()
	if err != nil {
		return false, err
	}
	result, ok := t.(bool)
	if !ok {
		return false, fmt.Errorf("parse error, expected bool at offset %v", j.c.Offset())
	}
	return result, nil
}

// Float returns the floating point number at the current position
func (j *JSONIterator) Float() (float64, error) {
	if !j.More() {
		return 0, fmt.Errorf("no more float available in current array or object at offset %v", j.c.Offset())
	}
	t, err := j.c.Token()
	if err != nil {
		return 0, err
	}
	result, ok := t.(float64)
	if !ok {
		return 0, fmt.Errorf("parse error, expected float at offset %v", j.c.Offset())
	}
	return result, nil
}

// String returns the string at the current position
func (j *JSONIterator) String() (string, error) {
	if !j.More() {
		return "", fmt.Errorf("no more string available in current array or object at offset %v", j.c.Offset())
	}
	t, err := j.c.Token()
	if err != nil {
		return "", err
	}
	result, ok := t.(string)
	if !ok || result == "" {
		return "", fmt.Errorf("parse error, expected string at offset %v", j.c.Offset())
	}
	return result, nil
}

func (j *JSONIterator) Range(f func(*Query) error) error {
	if err := j.Enter(); err != nil {
		return err
	}
	for j.More() {
		q, err := j.Next()
		if err != nil {
			return err
		}
		if err = f(q); err != nil {
			return err
		}
	}
	if err := j.Leave(); err != nil {
		return err
	}
	return nil
}

// Next returns the next element
func (j *JSONIterator) Next() (*Query, error) {
	return j.Query()
}

// Close closes the underlying ReaderCloser stream
func (j *JSONIterator) Close() error {
	return j.c.Close()
}

// PrintRaw copies the raw data to standard out
func (j *JSONIterator) PrintRaw() error {
	return j.c.PrintRaw()
}

// PrintPretty prints the raw data with indention to standard out
func (j *JSONIterator) PrintPretty() error {
	return j.c.PrintPretty()
}

// JSONMapIterator is a specialized JSONIterator. It returns the attribute name
// and content in an array. Use it to iterate JSON objects.
type JSONMapIterator struct {
	*JSONIterator
}

// NewJSONMapIterator creates an initialized JSONMapIterator
func NewJSONMapIterator(stream io.ReadCloser) *JSONMapIterator {
	base := NewJSONIterator(stream)
	return &JSONMapIterator{JSONIterator: base}
}

// Next returns the name and the content of the current attribute.
func (j *JSONMapIterator) Next() (*Query, error) {
	_, err := j.String()
	if err != nil {
		return nil, err
	}
	return j.Query()
}

func (j *JSONMapIterator) Range(f func(*Query) error) error {
	if err := j.Enter(); err != nil {
		return err
	}
	for j.More() {
		obj, err := j.Next()
		if err != nil {
			return err
		}
		if err = f(obj); err != nil {
			return err
		}
	}
	if err := j.Leave(); err != nil {
		return err
	}
	return nil
}

type JSONFilterIterator struct {
	*JSONIterator
	Filter map[string]bool
}

func NewJSONFilterIterator(filter []string, stream io.ReadCloser) *JSONFilterIterator {
	base := NewJSONIterator(stream)
	fm := make(map[string]bool, len(filter))
	for _, f := range filter {
		fm[f] = true
	}
	return &JSONFilterIterator{
		JSONIterator: base,
		Filter:       fm,
	}
}

func (j *JSONFilterIterator) Next() (*Query, error) {
	result := make(map[string]interface{}, len(j.Filter))
	for j.More() {
		id, err := j.String()
		if err != nil {
			return nil, err
		}
		if j.Filter[id] {
			obj, err := j.Interface()
			if err != nil {
				return nil, err
			}
			result[id] = obj
			continue
		}
		if err = j.Skip(); err != nil {
			return nil, err
		}
	}
	return NewQuery(result), nil
}

func (j *JSONFilterIterator) Range(f func(*Query) error) error {
	if err := j.Enter(); err != nil {
		return err
	}
	for j.More() {
		obj, err := j.Next()
		if err != nil {
			return err
		}
		if err = f(obj); err != nil {
			return err
		}
	}
	if err := j.Leave(); err != nil {
		return err
	}
	return nil
}
