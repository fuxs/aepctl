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

import "time"

// ObjectCacheData stores the objects
type ObjectCacheData struct {
	Map     []interface{}
	Expires time.Time
}

// Valid in checks if the token is still valid for passed duration
func (c *ObjectCacheData) Valid() bool {
	return time.Now().Before(c.Expires)
}

// ObjectCache is an object cache
type ObjectCache struct {
	exec    func() ([]interface{}, error)
	current *ObjectCacheData
	cache   *JSONCache
}

// NewObjectCache creates an initialized ObjectCache object
func NewObjectCache(f func() ([]interface{}, error)) *ObjectCache {
	return &ObjectCache{
		exec: f,
	}
}

// Load loads the cache
func (c *ObjectCache) Load() (*ObjectCacheData, error) {
	if c.current != nil {
		return c.current, nil
	}
	if c.cache != nil {
		if cd, _ := c.loadCacheData(); cd != nil {
			if cd.Valid() {
				return cd, nil
			}
		}
	}
	m, err := c.exec()
	if err != nil {
		return nil, err
	}
	result := &ObjectCacheData{
		Map:     m,
		Expires: time.Now().Add(time.Hour * 24),
	}
	if c.cache != nil {
		_ = c.cache.Save(result) // ignore error
	}
	c.current = result
	return result, nil
}

func (c *ObjectCache) loadCacheData() (*ObjectCacheData, error) {
	data := &ObjectCacheData{}
	if err := c.cache.Load(data); err != nil {
		return nil, err
	}
	return data, nil
}

// Mapper converts the stored objects to a Mapper. Parameters key and value are Query paths
// to the desired object attributes.
func (c *ObjectCache) Mapper(key, value []string) Mapper {
	data, err := c.Load()
	if err != nil {
		return nil
	}
	result := make(Mapper, len(data.Map))
	for _, v := range data.Map {
		q := NewQuery(v)
		result[q.Str(key...)] = q.Str(value...)
	}
	return result
}
