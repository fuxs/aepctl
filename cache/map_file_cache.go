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
	"time"

	"github.com/fuxs/aepctl/util"
)

// EatByMap stores data with an expiry date
type EatByMap struct {
	Map     util.Mapper
	Expires time.Time
}

// ValidIn checks if the token is still valid for passed duration
func (c *EatByMap) ValidIn(d time.Duration) bool {
	return time.Now().Add(d).Before(c.Expires)
}

//Valid checks if the token is still valid
func (c *EatByMap) Valid() bool {
	return time.Now().Before(c.Expires)
}

// MapFileCache is a file cache for maps
type MapFileCache struct {
	API      APICall
	Duration time.Duration
	trans    *TransformMap
	cache    *util.JSONFile
	cached   *EatByMap
}

// NewMapFileCache creates an intialized MapFileCache object
func NewMapFileCache(apiCall APICall, trans *TransformMap, d time.Duration, file util.Path) *MapFileCache {
	return &MapFileCache{
		API:      apiCall,
		Duration: d,
		trans:    trans,
		cache:    util.NewJSONFile(file),
	}
}

// DeleteE deletes the corresponding file
func (c *MapFileCache) Delete() {
	c.cached = nil
	_ = c.cache.Delete()
}

// DeleteE deletes the corresponding file
func (c *MapFileCache) DeleteE() error {
	c.cached = nil
	return c.cache.Delete()
}

// Load loads the cache
func (c *MapFileCache) Load() error {
	// already loaded?
	if c.cached != nil && c.cached.Valid() {
		return nil
	}
	// load from disk
	eb := &EatByMap{}
	if err := c.cache.Load(&eb); err == nil {
		if eb.Valid() {
			c.cached = eb
			return nil
		}
	}
	// get data from server
	obj, err := c.API.Call()
	if err != nil {
		return err
	}
	result := &EatByMap{
		Map:     c.trans.Transform(obj),
		Expires: time.Now().Add(c.Duration),
	}
	c.cached = result
	// save result
	_ = c.cache.Save(result)
	return nil
}

// Lookup returns either the related value or the key itself
func (c *MapFileCache) Lookup(key string) string {
	if err := c.Load(); err != nil {
		return key
	}
	return c.cached.Map.Lookup(key)
}

// LookupE returns either the related value or an error
func (c *MapFileCache) LookupE(key string) (string, error) {
	if err := c.Load(); err != nil {
		return "", err
	}
	return c.cached.Map.Lookup(key), nil
}

// Keys returns all keys from the map
func (c *MapFileCache) Keys() []string {
	if err := c.Load(); err != nil {
		return []string{}
	}
	return c.cached.Map.Keys()
}

// Values returns all values from the map
func (c *MapFileCache) Values() []string {
	if err := c.Load(); err != nil {
		return []string{}
	}
	return c.cached.Map.Values()
}
