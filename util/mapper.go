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
	"strings"
)

// Mapper is a map of strings
type Mapper map[string]string

// Get returns the value for the passed key. If no value can be found the key will be returned.
func (m Mapper) Lookup(key string) string {
	if result, ok := m[key]; ok {
		return result
	}
	return key
}

func (m Mapper) LookupE(key string) (string, error) {
	if result, ok := m[key]; ok {
		return result, nil
	}
	return "", fmt.Errorf("Could not find key %v", key)
}

// GetL normalizes the key to lower case before getting the value. If no value can be found the original key will be returned.
func (m Mapper) GetL(key string) string {
	nkey := strings.ToLower(key)
	if result, ok := m[nkey]; ok {
		return result
	}
	return key
}

// Keys returns a list of keys
func (m Mapper) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns a list of values
func (m Mapper) Values() []string {
	values := make([]string, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Invert creats a new mapper with inverted key value relation. Key -> Value becomes Value -> Key.
func (m Mapper) Invert() Mapper {
	result := make(Mapper, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}
