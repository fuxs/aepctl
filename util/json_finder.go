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
)

type pathAndFunc struct {
	path []string
	f    func(JSONResponse) error
}

// StackInt is a simple int stack (lifo)
type stack [][]*pathAndFunc

// Push pushes the integer i to the stack
func (s *stack) Push(i []*pathAndFunc) {
	*s = append(*s, i)
}

// Pop removes the last element from the stack and returns it
func (s *stack) Pop() ([]*pathAndFunc, bool) {
	l := len(*s)
	if l == 0 {
		return nil, false
	}
	result := (*s)[l-1]
	*s = (*s)[:l-1]
	return result, true
}

type Prefixes struct {
	history stack
	current []*pathAndFunc
}

func NewPrefixes() *Prefixes {
	return &Prefixes{current: make([]*pathAndFunc, 0, 2), history: make(stack, 0, 8)}
}

func (p *Prefixes) Add(f func(JSONResponse) error, path ...string) {
	p.current = append(p.current, &pathAndFunc{path: path, f: f})
}

func (p *Prefixes) Matches(pos int, elem string) (bool, func(JSONResponse) error) {
	next := make([]*pathAndFunc, 0, len(p.current))
	found := false
	for _, n := range p.current {
		if n.path[pos] == elem {
			if pos == len(n.path)-1 {
				return true, n.f
			}
			found = true
			next = append(next, n)
		}
	}
	if found {
		p.history.Push(p.current)
		p.current = next
	}
	return found, nil
}

func (p *Prefixes) Pop() {
	if current, ok := p.history.Pop(); ok {
		p.current = current
	}
}

type JSONFinder struct {
	i        JSONResponse
	prefixes *Prefixes
	depth    int
}

func NewJSONFinder() *JSONFinder {
	return &JSONFinder{prefixes: NewPrefixes()}
}

func (j *JSONFinder) SetIterator(i JSONResponse) {
	j.i = i
}

func (j *JSONFinder) Add(f func(JSONResponse) error, path ...string) {
	j.prefixes.Add(f, path...)
}

// Path skips all element along the path
func (j *JSONFinder) Run() error {
	i := j.i
	if err := i.EnterObject(); err != nil {
		return err
	}
	j.depth++

	for j.depth > 0 {
		// get the next token
		t, err := i.Token()
		if err != nil {
			return err
		}
		// if delimiter, only } can be accepted
		d, ok := t.(json.Delim)
		if ok {
			if d != '}' {
				return fmt.Errorf("unexpected token %v", d)
			}
			j.depth--
			if j.depth > 0 {
				j.prefixes.Pop()
			}
			continue
		}
		// it must be a string
		id, ok := t.(string)
		if !ok || id == "" {
			return fmt.Errorf("parse error, expected string at offset %v", i.Offset())
		}
		found, f := j.prefixes.Matches(j.depth-1, id)
		if found {
			if f != nil {
				if err = f(i); err != nil {
					return err
				}
			} else {
				if err = i.EnterObject(); err != nil {
					return err
				}
				j.depth++
			}
		} else {
			if err = i.Skip(); err != nil {
				return err
			}
		}
	}
	return nil
}
