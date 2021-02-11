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

import "gopkg.in/yaml.v3"

// YAMLQuery supports the querying and editing of YAML documents without loosing
// the formating
type YAMLQuery struct {
	node *yaml.Node
}

// NewYAMLQuery creates an initialized YAMLQuery object
func NewYAMLQuery(node *yaml.Node) *YAMLQuery {
	return &YAMLQuery{node: node}
}

// Path queries the passed path. Each path element must match an entry in the
// current map, otherwise an empty object will be returned.
func (q *YAMLQuery) Path(path ...string) *YAMLQuery {
	if len(path) == 0 {
		return q
	}
	r := q.First()
	if !r.IsMap() {
		return &YAMLQuery{}
	}
	cur := r.node
	found := true
	for _, p := range path {
		if cur.Kind != yaml.MappingNode {
			return &YAMLQuery{}
		}
		found = false
		for i := 0; i < len(cur.Content); i = i + 2 {
			key := cur.Content[i]
			if key.Value == p {
				cur = cur.Content[i+1]
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	if found {
		return &YAMLQuery{node: cur}
	}
	return &YAMLQuery{}
}

// IsNil tests if the current object is nil
func (q *YAMLQuery) IsNil() bool {
	return q.node == nil
}

// String returns the current object as string
func (q *YAMLQuery) String() string {
	if q.node == nil {
		return ""
	}
	return q.node.Value
}

// Str returns the string of the referenced path
func (q *YAMLQuery) Str(path ...string) string {
	return q.Path(path...).String()
}

// Set sets the value of the current object
func (q *YAMLQuery) Set(value string) {
	if q.node != nil {
		q.node.Value = value
		q.node.Tag = "!!str"
	}
}

// First skips the document node if necessary and returns the first real object
func (q *YAMLQuery) First() *YAMLQuery {
	if q.node != nil && len(q.node.Content) > 0 {
		if q.node.Kind == yaml.DocumentNode {
			return &YAMLQuery{node: q.node.Content[0]}
		}
	}
	return q
}

// IsMap returns true if this object is a map
func (q *YAMLQuery) IsMap() bool {
	return q.node != nil && q.node.Kind == yaml.MappingNode
}

// IsDocument returns true if this object is a document
func (q *YAMLQuery) IsDocument() bool {
	return q.node != nil && q.node.Kind == yaml.DocumentNode
}

// SetMap sets the key-value pair
func (q *YAMLQuery) SetMap(key, value string) {
	r := q.First()
	if r.IsMap() {
		c := r.node.Content
		l := len(c)
		for i := 0; i < l; i = i + 2 {
			k := c[i]
			if k.Value == key {
				c[i+1].Value = value
				return
			}
		}
		// append new entry
		//row := q.LastLine() + 1
		node := &yaml.Node{
			Kind: yaml.ScalarNode,
			Tag:  "!!str",
			//Column: 1,
			//Line:   row,
			Value: key,
		}
		c = append(c, node)
		node = &yaml.Node{
			Kind: yaml.ScalarNode,
			Tag:  "!!str",
			//Column: 1,
			//Line:   row + 1,
			Value: value,
		}
		c = append(c, node)
		r.node.Content = c
	}

}
