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
	"time"
)

// StringCacheData stores a list of strings with expiry date
type StringCacheData struct {
	List    []string
	Expires time.Time
}

// ValidIn in checks if the list is still valid for passed duration
func (c *StringCacheData) ValidIn(d time.Duration) bool {
	return time.Now().Add(d).Before(c.Expires)
}

// StringCache is a string cache
type StringCache struct {
	cache      *JSONCache
	exec       func() ([]string, error)
	cacheRead  func() bool
	cacheWrite func() bool
}

// NewStringCache creates an initialized StringCache object
func NewStringCache(cacheRead, cacheWrite func() bool, exec func() ([]string, error), path func() []string) *StringCache {
	return &StringCache{
		cache:      NewJSONCache(path),
		exec:       exec,
		cacheRead:  cacheRead,
		cacheWrite: cacheWrite,
	}
}

func (l *StringCache) loadCache() (*StringCacheData, error) {
	list := &StringCacheData{}
	if err := l.cache.Load(list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetList returns all stored values
func (l *StringCache) GetList() ([]string, error) {
	ss, err := l.Load()
	if err != nil {
		return []string{}, err
	}
	return ss.List, nil
}

// Load loads the cache
func (l *StringCache) Load() (*StringCacheData, error) {
	if l.cacheRead() && l.cache != nil {
		if sl, _ := l.loadCache(); sl != nil {
			if sl.ValidIn(time.Minute) {
				return sl, nil
			}
		}
	}
	list, err := l.exec()
	if err != nil {
		return nil, err
	}
	result := &StringCacheData{
		List:    list,
		Expires: time.Now().Add(time.Hour * 24),
	}
	if l.cacheWrite() && l.cache != nil {
		_ = l.cache.Save(result) // ignore error
	}
	return result, nil
}
