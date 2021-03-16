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

// PathProvider is the interface for objects providing dynamic file path
// concatenation
type PathProvider interface {
	Path(...string) string
}

// Path is the interface for objects providing a file path
type Path interface {
	Path() string
}

// LazyPath generates a file path on demand
type LazyPath struct {
	pp     PathProvider
	path   []string
	result string
}

// NewLazyPath creates an initialzed LazyPath object
func NewLazyPath(pp PathProvider, path ...string) *LazyPath {
	return &LazyPath{
		pp:   pp,
		path: path,
	}
}

// Path generates the path
func (l *LazyPath) Path() string {
	if l.result == "" {
		l.result = l.pp.Path(l.path...)
	}
	return l.result
}
