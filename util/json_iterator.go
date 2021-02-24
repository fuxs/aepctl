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
	"fmt"
	"io"
	"os"
)

// JSONResponse is the interface for streaming JSON objects
type JSONResponse interface {
	Close() error
	Delim() (json.Delim, error)
	EnterArray() error
	EnterObject() error
	More() bool
	Next() (interface{}, error)
	Obj() (map[string]interface{}, error)
	Path(...string) error
	//	Stream() io.ReadCloser
	PrintRaw() error
	PrintPretty() error
}

// JSONIterator implements the JSONResponse interface
type JSONIterator struct {
	dec    *json.Decoder
	stream io.ReadCloser
}

// NewJSONIterator creates an initialized JSONIterator object
func NewJSONIterator(stream io.ReadCloser) (*JSONIterator, error) {
	return &JSONIterator{dec: json.NewDecoder(stream), stream: stream}, nil
}

// More checks if there is another element in the current object or array
func (j *JSONIterator) More() bool {
	return j.dec != nil && j.dec.More()
}

// Skip skips the next element like for example a string, object or array
func (j *JSONIterator) Skip() error {
	t, err := j.dec.Token()
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
		for t, err = j.dec.Token(); err != nil; t, err = j.dec.Token() {
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
		}
		return err
	}
	// must be '['
	for t, err = j.dec.Token(); err != nil; t, err = j.dec.Token() {
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
	}
	return err
}

// EnterArray moves forward to the first element in the array
func (j *JSONIterator) EnterArray() error {
	d, err := j.Delim()
	if err != nil {
		return err
	}
	if d != '[' {
		return fmt.Errorf("Expecting [ found %v at offset %v", d, j.dec.InputOffset())
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
		return fmt.Errorf("Expected { at %v but found %v", j.dec.InputOffset(), d)
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
			if id == p {
				found = true
				break
			}
			if err = j.Skip(); err != nil {
				return err
			}
		}
		if !found {
			return fmt.Errorf("Could not find id %v but reached end at offset %v", p, j.dec.InputOffset())
		}
	}
	return nil
}

// Delim returns the delimiter like [,],{ or} at the current position
func (j *JSONIterator) Delim() (json.Delim, error) {
	t, err := j.dec.Token()
	if err != nil {
		return 0, err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return 0, fmt.Errorf("Parse error, expected delimiter at offset %v", j.dec.InputOffset())
	}
	return delim, nil
}

// Obj returns the object at the current position
func (j *JSONIterator) Obj() (map[string]interface{}, error) {
	if !j.More() {
		return nil, fmt.Errorf("No more object available in current array or object at offset %v", j.dec.InputOffset())
	}
	var obj map[string]interface{}
	if err := j.dec.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// Interface returns any element at the current position
func (j *JSONIterator) Interface() (interface{}, error) {
	if !j.More() {
		return nil, fmt.Errorf("No more element available in current array or object at offset %v", j.dec.InputOffset())
	}
	var obj interface{}
	if err := j.dec.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// Bool returns the boolean at the current position
func (j *JSONIterator) Bool() (bool, error) {
	if !j.More() {
		return false, fmt.Errorf("No more bool available in current array or object at offset %v", j.dec.InputOffset())
	}
	t, err := j.dec.Token()
	if err != nil {
		return false, err
	}
	result, ok := t.(bool)
	if !ok {
		return false, fmt.Errorf("Parse error, expected bool at offset %v", j.dec.InputOffset())
	}
	return result, nil
}

// Float returns the floating point number at the current position
func (j *JSONIterator) Float() (float64, error) {
	if !j.More() {
		return 0, fmt.Errorf("No more float available in current array or object at offset %v", j.dec.InputOffset())
	}
	t, err := j.dec.Token()
	if err != nil {
		return 0, err
	}
	result, ok := t.(float64)
	if !ok {
		return 0, fmt.Errorf("Parse error, expected float at offset %v", j.dec.InputOffset())
	}
	return result, nil
}

// String returns the string at the current position
func (j *JSONIterator) String() (string, error) {
	if !j.More() {
		return "", fmt.Errorf("No more string available in current array or object at offset %v", j.dec.InputOffset())
	}
	t, err := j.dec.Token()
	if err != nil {
		return "", err
	}
	result, ok := t.(string)
	if !ok || result == "" {
		return "", fmt.Errorf("Parse error, expected string at offset %v", j.dec.InputOffset())
	}
	return result, nil
}

// Next returns the next element
func (j *JSONIterator) Next() (interface{}, error) {
	return j.Interface()
}

// Close closes the underlying ReaderCloser stream
func (j *JSONIterator) Close() error {
	return j.stream.Close()
}

// PrintRaw copies the raw data to standard out
func (j *JSONIterator) PrintRaw() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	_, err := io.Copy(bout, j.stream)
	return err
}

// PrintPretty prints the raw data with indention to standard out
func (j *JSONIterator) PrintPretty() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	return JSONPrintPretty(j.dec, bout)
}

// JSONMapIterator is a specialized JSONIterator. It returns the attribute name
// and content in an array. Use it to iterate JSON objects.
type JSONMapIterator struct {
	*JSONIterator
}

// NewJSONMapIterator creates an initialized JSONMapIterator
func NewJSONMapIterator(stream io.ReadCloser) (*JSONMapIterator, error) {
	base, err := NewJSONIterator(stream)
	if err != nil {
		return nil, err
	}
	return &JSONMapIterator{JSONIterator: base}, nil
}

// Next returns the name and the content of the current attribute.
func (j *JSONMapIterator) Next() (interface{}, error) {
	id, err := j.String()
	if err != nil {
		return nil, err
	}
	obj, err := j.Interface()
	if err != nil {
		return nil, err
	}
	return []interface{}{id, obj}, nil
}
