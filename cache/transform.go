/*
Package cache consists of all caching relted functions and data structures.

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
package cache

import "github.com/fuxs/aepctl/util"

// TransformMap transforms a JSON response to a mapper
type TransformMap struct {
	Path  []string
	Key   []string
	Value []string
}

// NewTransformMap creates an initialzed TransformMap object
func NewTransformMap(p ...string) *TransformMap {
	return &TransformMap{Path: p}
}

// P provides the path to the rangable list inside the JSON file
func (c *TransformMap) P(p ...string) *TransformMap {
	c.Path = p
	return c
}

// K provides the path to the key
func (c *TransformMap) K(k ...string) *TransformMap {
	c.Key = k
	return c
}

// V provides the path to the value
func (c *TransformMap) V(v ...string) *TransformMap {
	c.Value = v
	return c
}

// Transform transforms the passed ojbect to a mapper
func (c TransformMap) Transform(obj interface{}) util.Mapper {
	m := make(util.Mapper)
	util.NewQuery(obj).Path(c.Path...).Range(func(q *util.Query) {
		m[q.Str(c.Key...)] = q.Str(c.Value...)
	})
	return m
}

// TransformList transforms a JSON response to a list
type TransformList struct {
	Path  []string
	Value []string
}

// NewTransformList creates an initialized TransformList object
func NewTransformList(p ...string) *TransformList {
	return &TransformList{Path: p}
}

// P provides the path to the rangable list inside the JSON file
func (c *TransformList) P(p ...string) *TransformList {
	c.Path = p
	return c
}

// V provides the path to the value
func (c *TransformList) V(v ...string) *TransformList {
	c.Value = v
	return c
}

// Transform transforms the passed ojbect to an array of strings
func (c TransformList) Transform(obj interface{}) []string {
	var m []string
	//m := make(util.Mapper, 0)
	util.NewQuery(obj).Path(c.Path...).Range(func(q *util.Query) {
		m = append(m, q.Str(c.Value...))
	})
	return m
}
