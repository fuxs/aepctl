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

// EatByList stores data with an expiry date
type EatByList struct {
	List    []string
	Expires time.Time
}

// ValidIn in checks if the token is still valid for passed duration
func (c *EatByList) ValidIn(d time.Duration) bool {
	return time.Now().Add(d).Before(c.Expires)
}

func (c *EatByList) Valid() bool {
	return time.Now().Before(c.Expires)
}

type ListFileCache struct {
	API      APICall
	Duration time.Duration
	trans    *TransformList
	cache    *util.JSONFile
	cached   *EatByList
}

func NewListFileCache(apiCall APICall, trans *TransformList, d time.Duration, file util.Path) *ListFileCache {
	return &ListFileCache{
		API:      apiCall,
		Duration: d,
		trans:    trans,
		cache:    util.NewJSONFile(file),
	}
}

// Delete deletes the corresponding file
func (c *ListFileCache) Delete() error {
	c.cached = nil
	return c.cache.Delete()
}

// Load loads the cache
func (c *ListFileCache) Load() error {
	// already loaded?
	if c.cached != nil && c.cached.Valid() {
		return nil
	}
	// load from disk
	eb := &EatByList{}
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
	result := &EatByList{
		List:    c.trans.Transform(obj),
		Expires: time.Now().Add(c.Duration),
	}
	c.cached = result
	// save result
	_ = c.cache.Save(result)
	return nil
}

func (c *ListFileCache) Values() []string {
	if err := c.Load(); err != nil {
		return []string{}
	}
	return c.cached.List
}
