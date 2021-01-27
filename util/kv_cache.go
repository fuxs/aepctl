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
	"fmt"
	"time"
)

// KVCacheData stores the data of the KVCache with an expiry date
type KVCacheData struct {
	Map     map[string]string
	Expires time.Time
}

// ValidIn in checks if the token is still valid for passed duration
func (c *KVCacheData) ValidIn(d time.Duration) bool {
	return time.Now().Add(d).Before(c.Expires)
}

// Get returns the value for the passed key. If no value could be found then the key will be returned.
func (c *KVCacheData) Get(key string) string {
	result, ok := c.Map[key]
	if !ok {
		return key
	}
	return result
}

// GetE returns the value for the passed key. If no value could be found then an error will be returned.
func (c *KVCacheData) GetE(key string) (string, error) {
	result, ok := c.Map[key]
	if !ok {
		return "", fmt.Errorf("No value for key %s", key)
	}
	return result, nil
}

// Keys returns a list of keys
func (c *KVCacheData) Keys() []string {
	keys := make([]string, len(c.Map))
	i := 0
	for k := range c.Map {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns a list of values
func (c *KVCacheData) Values() []string {
	values := make([]string, len(c.Map))
	i := 0
	for _, v := range c.Map {
		values[i] = v
		i++
	}
	return values
}

// KVCache is a key value cache for strings with optional persistence
type KVCache struct {
	exec       func() (map[string]string, error)
	cacheRead  func() bool
	cacheWrite func() bool
	current    *KVCacheData
	cache      *JSONCache
}

func (l *KVCache) loadCache() (*KVCacheData, error) {
	list := &KVCacheData{}
	if err := l.cache.Load(list); err != nil {
		return nil, err
	}
	return list, nil
}

func (l *KVCache) save(list *KVCacheData) error {
	return l.cache.Save(list)
}

// NewKVCache creates an initialized KVCache object
func NewKVCache(cacheRead, cacheWrite func() bool, exec func() (map[string]string, error), path func() []string) *KVCache {
	if path == nil {
		return &KVCache{
			exec:       exec,
			cacheRead:  cacheRead,
			cacheWrite: cacheWrite,
		}
	}
	return &KVCache{
		cache:      NewJSONCache(path),
		exec:       exec,
		cacheRead:  cacheRead,
		cacheWrite: cacheWrite,
	}
}

// Invalidate invalidates the cache
func (l *KVCache) Invalidate() {
	if l.cache != nil {
		l.cache.Invalidate()
	}
	l.current = nil
}

// Remove removes the string with the passed key
func (l *KVCache) Remove(key string) {
	store, err := l.Load()
	if err != nil {
		return
	}
	delete(store.Map, key)
	if l.cacheWrite() {
		_ = l.save(store) // ignore error
	}
}

// GetValue returns the value for the passed key. If no value could be retrieved then the key will be returned.
func (l *KVCache) GetValue(key string) string {
	store, err := l.Load()
	if err != nil {
		return key
	}
	return store.Get(key)
}

// GetValueE returns the value for the passed key. If no value could be retrieved then an error will be returned.
func (l *KVCache) GetValueE(key string) (string, error) {
	store, err := l.Load()
	if err != nil {
		return "", err
	}
	return store.GetE(key)
}

// Keys returns a list of keys
func (l *KVCache) Keys() ([]string, error) {
	store, err := l.Load()
	if err != nil {
		return []string{}, err
	}
	return store.Keys(), nil
}

// Values returns a list of values
func (l *KVCache) Values() ([]string, error) {
	store, err := l.Load()
	if err != nil {
		return []string{}, err
	}
	return store.Values(), nil
}

// MapValues calls the passed function for each value and stores the result
func (l *KVCache) MapValues(m func(string) string) {
	store, err := l.Load()
	if err != nil {
		return
	}
	for k, v := range store.Map {
		store.Map[k] = m(v)
	}
}

// Load loads the cache
func (l *KVCache) Load() (*KVCacheData, error) {
	if l.current != nil {
		return l.current, nil
	}
	if l.cache != nil && l.cacheRead() {
		if sl, _ := l.loadCache(); sl != nil {
			if sl.ValidIn(time.Minute) {
				return sl, nil
			}
		}
	}
	m, err := l.exec()
	if err != nil {
		return nil, err
	}
	result := &KVCacheData{
		Map:     m,
		Expires: time.Now().Add(time.Hour * 24),
	}
	if l.cache != nil && l.cacheWrite() {
		_ = l.save(result) // ignore error
	}
	l.current = result
	return result, nil
}
