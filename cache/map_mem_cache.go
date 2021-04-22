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

import (
	"github.com/fuxs/aepctl/util"
)

// MapMemCache stores a mapper object im memory
type MapMemCache struct {
	API    APICall
	trans  *TransformMap
	cached util.Mapper
}

// NewMapMemCache creates an initialized MapMemCache object
func NewMapMemCache(apiCall APICall, trans *TransformMap) *MapMemCache {
	return &MapMemCache{API: apiCall, trans: trans}
}

// Delete just empties the cache
func (c *MapMemCache) Delete() error {
	c.cached = nil
	return nil
}

// Load loads the cache
func (c *MapMemCache) Load() error {
	if c.cached == nil {
		obj, err := c.API.Call()
		if err != nil {
			return err
		}
		c.cached = c.trans.Transform(obj)
	}
	return nil
}

// LookupE returns either the related value or the key
func (c *MapMemCache) Lookup(key string) string {
	if err := c.Load(); err != nil {
		return key
	}
	return c.cached.Lookup(key)
}

// LookupE returns either the related value or an error
func (c *MapMemCache) LookupE(key string) (string, error) {
	if err := c.Load(); err != nil {
		return "", err
	}
	return c.cached.Lookup(key), nil
}

// Keys returns all keys from the map
func (c *MapMemCache) Keys() []string {
	if err := c.Load(); err != nil {
		return []string{}
	}
	return c.cached.Keys()
}

// Values returns all values from the map
func (c *MapMemCache) Values() []string {
	if err := c.Load(); err != nil {
		return []string{}
	}
	return c.cached.Values()
}

func (c *MapMemCache) Mapper() util.Mapper {
	if err := c.Load(); err != nil {
		return nil
	}
	return c.cached
}
