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
	"io/ioutil"
	"os"
)

type MultiFileReader struct {
	Files   []string
	Current int
}

func (m *MultiFileReader) Read() ([]byte, error) {
	l := len(m.Files)
	if m.Current == 0 && l == 0 {
		m.Current++
		return ioutil.ReadAll(os.Stdin)
	}
	if m.Current >= l {
		return nil, nil
	}
	path := m.Files[m.Current]
	m.Current++
	file := os.Stdin
	if path != "-" {
		var err error
		if file, err = os.Open(path); err != nil {
			return nil, err
		}
	}
	return ioutil.ReadAll(file)
}

func (m *MultiFileReader) ReadAll(f func(data []byte) error) error {
	if f == nil {
		return nil
	}
	data, err := m.Read()
	// data != nil means no more data
	for err == nil && data != nil {
		if err = f(data); err == nil {
			data, err = m.Read()
		}
	}
	return err
}
