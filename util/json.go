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
	"strings"
)

const (
	empty int = iota
	mapState
	arrayState
)

func JSONPrintPrettyln(dec *json.Decoder, out io.Writer) error {
	if err := JSONPrintPretty(dec, out); err != nil {
		return err
	}
	fmt.Fprintln(out)
	return nil
}

// JSONPrintPretty streams the dec JSON content in an indented format to out
func JSONPrintPretty(dec *json.Decoder, out io.Writer) error {
	var (
		expectID    bool
		token, peek interface{}
	)

	state := empty
	depth := 0
	curIndent := ""
	indent := func() error {
		_, err := out.Write([]byte(curIndent))
		return err
	}
	newLine := func() error {
		var err error
		switch state {
		case arrayState:
			peek, err = dec.Token()
			if err != nil {
				return err
			}
			if p, ok := peek.(json.Delim); !ok || p != ']' {
				fmt.Fprintln(out, ",")
			} else {
				depth--
				curIndent = strings.Repeat("  ", depth)
				fmt.Fprintln(out)
			}
		case mapState:
			peek, err = dec.Token()
			if err != nil {
				return err
			}
			if p, ok := peek.(json.Delim); !ok || p != '}' {
				fmt.Fprintln(out, ",")
				expectID = true
			} else {
				depth--
				curIndent = strings.Repeat("  ", depth)
				fmt.Fprintln(out)
			}
		}
		return indent()
	}
	var (
		err error
		ok  bool
	)
	ss := new(StackInt)
	for {
		if peek != nil {
			token = peek
			peek = nil
		} else {
			token, err = dec.Token()
			if err != nil {
				return err
			}
		}
		// an expected id can only be a string
		if expectID {
			id, ok := token.(string)
			if !ok {
				return fmt.Errorf("expected id at position %v", dec.InputOffset())
			}
			if _, err := out.Write([]byte(JSONString(id))); err != nil {
				return err
			}
			if _, err := out.Write([]byte(": ")); err != nil {
				return err
			}
			expectID = false
			continue
		}
		// no id
		switch t := token.(type) {
		case json.Delim:
			switch t {
			case '[':
				// check for empty array
				peek, err = dec.Token()
				if err != nil {
					return err
				}
				if p, ok := peek.(json.Delim); ok && p == ']' {
					peek = nil
					fmt.Fprint(out, "[]")
					if err = newLine(); err != nil {
						return err
					}
					if state == empty {
						return nil
					}
					break
				}
				// not empty
				fmt.Fprintln(out, "[")
				ss.Push(state)
				state = arrayState
				depth++
				curIndent = strings.Repeat("  ", depth)
				if err = indent(); err != nil {
					return err
				}
			case ']':
				if state, ok = ss.Pop(); !ok {
					return fmt.Errorf("unbalanced ] at postion %v", dec.InputOffset())
				}
				if _, err = out.Write([]byte("]")); err != nil {
					return err
				}
				if err = newLine(); err != nil {
					return err
				}
				if state == empty {
					return nil
				}

			case '{':
				// check for empty object
				peek, err = dec.Token()
				if err != nil {
					return err
				}
				if p, ok := peek.(json.Delim); ok && p == '}' {
					peek = nil
					fmt.Fprint(out, "{}")
					if err = newLine(); err != nil {
						return err
					}
					if state == empty {
						return nil
					}
					break
				}
				// not empty
				fmt.Fprintln(out, "{")
				ss.Push(state)
				state = mapState
				expectID = true

				depth++
				curIndent = strings.Repeat("  ", depth)
				if err = indent(); err != nil {
					return err
				}
			case '}':
				if state, ok = ss.Pop(); !ok {
					return fmt.Errorf("unbalanced ] at postion %v", dec.InputOffset())
				}
				if _, err = out.Write([]byte("}")); err != nil {
					return err
				}
				if err = newLine(); err != nil {
					return err
				}
				if state == empty {
					return nil
				}
			}
		case string:
			if _, err = out.Write([]byte(JSONString(t))); err != nil {
				return err
			}
			if err = newLine(); err != nil {
				return err
			}
		case float32:
			if _, err = out.Write([]byte(JSONFloat(float64(t), 32))); err != nil {
				return err
			}
			if err = newLine(); err != nil {
				return err
			}
		case float64:
			if _, err = out.Write([]byte(JSONFloat(t, 64))); err != nil {
				return err
			}
			if err = newLine(); err != nil {
				return err
			}
		case bool:
			if t {
				if _, err = out.Write([]byte("true")); err != nil {
					return err
				}
			} else {
				if _, err = out.Write([]byte("false")); err != nil {
					return err
				}
			}
			if err = newLine(); err != nil {
				return err
			}
		default:
			if err = newLine(); err != nil {
				return err
			}
		}
	}
}
