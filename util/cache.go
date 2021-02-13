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
	"io/ioutil"
	"os"
	"path/filepath"
)

// LazyPath builds a path on the first call
type LazyPath struct {
	exec func() []string
	full []string
	path []string
	name string
}

// NewLazyPath returns an initialized LazyPath object
func NewLazyPath(exec func() []string) *LazyPath {
	return &LazyPath{exec: exec}
}

func (lp *LazyPath) d() *LazyPath {
	if lp.full == nil {
		lp.full = lp.exec()
		l := len(lp.full)
		if l > 1 {
			lp.path = lp.full[0 : l-1]
			lp.name = lp.full[l-1]
		}
		if l == 1 {
			lp.path = []string{}
			lp.name = lp.full[0]
		}
	}
	return lp
}

// Full returns the full path
func (lp *LazyPath) Full() []string {
	return lp.d().full
}

// Path returns the path component
func (lp *LazyPath) Path() []string {
	return lp.d().path
}

// Name returns the name of the file
func (lp *LazyPath) Name() string {
	return lp.d().name
}

// JSONCache stores json objects in a file
type JSONCache struct {
	lp *LazyPath
}

// NewJSONCache returns an initialzed JSONCache object
func NewJSONCache(path func() []string) *JSONCache {
	return &JSONCache{
		lp: NewLazyPath(path),
	}
}

// Invalidate invalidates the cache by deleting the related file
func (jc *JSONCache) Invalidate() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	p := filepath.Join(append([]string{home, ".aepctl", "cache"}, jc.lp.Full()...)...)
	os.Remove(p)
}

// Save stores the passed object in json format to a file
func (jc *JSONCache) Save(obj interface{}) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	p := filepath.Join(append([]string{home, ".aepctl", "cache"}, jc.lp.Path()...)...)
	if _, err = os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, 0700); err != nil {
			return err
		}
	}
	var data []byte
	data, err = json.Marshal(obj)
	if err != nil {
		return err
	}
	p = filepath.Join(p, jc.lp.Name())
	return ioutil.WriteFile(p, data, 0700)
}

// Load loads the json file into the passed object
func (jc *JSONCache) Load(obj interface{}) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	p := filepath.Join(append([]string{home, ".aepctl", "cache"}, jc.lp.Full()...)...)
	var data []byte
	data, err = ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &obj)
}
