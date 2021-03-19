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

// JSONFile stores json objects in a file
type JSONFile struct {
	p Path
}

// NewJSONFile returns an initialzed JSONCache object
func NewJSONFile(path Path) *JSONFile {
	return &JSONFile{p: path}
}

// Delete deletes the related file
func (jc *JSONFile) Delete() error {
	return os.Remove(jc.p.Path())
}

// Save stores the passed object in json format to a file
func (jc *JSONFile) Save(obj interface{}) error {
	p := filepath.Dir(jc.p.Path())
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, 0700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jc.p.Path(), data, 0700)
}

// Exists returns true if the related file exists
func (jc *JSONFile) Exists() bool {
	_, err := os.Stat(jc.p.Path())
	return err != nil
}

// Load loads the json file into the passed object
func (jc *JSONFile) Load(obj interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(jc.p.Path())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &obj)
}
