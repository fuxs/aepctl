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
	"os"
	"path/filepath"
	"sort"
)

// Dir helps to browse directories
type Dir struct {
	Path    string
	parent  string
	files   []os.FileInfo
	current []os.FileInfo
}

// NewDir creates an initialized Dir object for the current directory
func NewDir() (*Dir, error) {
	return NewDirPath(".")
}

// NewDirPath creates an initialized Dir object for the passed path
func NewDirPath(dir string) (*Dir, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return &Dir{
		Path:   abs,
		parent: filepath.Dir(abs),
		files:  files,
	}, nil
}

func (d *Dir) all(hidden bool) []os.FileInfo {
	if hidden {
		// copy the slice
		return append([]os.FileInfo(nil), d.files...)
	}
	// remove hidden files
	result := make([]os.FileInfo, 0, len(d.files))
	for _, f := range d.files {
		if f.Name()[0:1] != "." {
			result = append(result, f)
		}
	}
	return result
}

// IsDir returns true if the ith child is a directory
func (d *Dir) IsDir(i int) bool {
	return d.current[i].IsDir()
}

// HasParent returns true if this Dir object is not root
func (d *Dir) HasParent() bool {
	return d.Path != "/"
}

// Parent returns the parent Dir object or nil if this Dir object is root
func (d *Dir) Parent() (*Dir, error) {
	if d.HasParent() {
		return NewDirPath(d.parent)
	}
	return nil, nil
}

// PathI returns the full path for the ith child
func (d *Dir) PathI(i int) string {
	return filepath.Join(d.Path, d.current[i].Name())
}

// ChildI returns the Dir object of the ith child
func (d *Dir) ChildI(i int) (*Dir, error) {
	return d.Child(d.current[i].Name())
}

// Child returns the Dir object for the child with the passed name
func (d *Dir) Child(name string) (*Dir, error) {
	return NewDirPath(filepath.Join(d.Path, name))
}

// SortedByName returns a list of FileInfo objects sorted by name
func (d *Dir) SortedByName(hidden, sep, asc bool) []os.FileInfo {
	files := d.all(hidden)
	sort.Slice(files, func(i, j int) bool {
		fi, fj := files[i], files[j]
		if sep {
			if fi.IsDir() {
				if !fj.IsDir() {
					return true
				}
			} else {
				if fj.IsDir() {
					return false
				}
			}
		}
		r := fi.Name() < fj.Name()
		if asc {
			return r
		}
		return !r
	})
	d.current = files
	return files
}

// SortedBySize returns a list of FileInfo objects sorted by size
func (d *Dir) SortedBySize(hidden, sep, asc bool) []os.FileInfo {
	files := d.all(hidden)
	sort.Slice(files, func(i, j int) bool {
		fi, fj := files[i], files[j]
		if sep {
			if fi.IsDir() {
				if !fj.IsDir() {
					return true
				}
			} else {
				if fj.IsDir() {
					return false
				}
			}
		}
		r := fi.Size() < fj.Size()
		if asc {
			return r
		}
		return !r
	})
	d.current = files
	return files
}

// SortedByModTime returns a list of FileInfo objects sorted by modification
// date
func (d *Dir) SortedByModTime(hidden, sep, asc bool) []os.FileInfo {
	files := d.all(hidden)
	sort.Slice(files, func(i, j int) bool {
		fi, fj := files[i], files[j]
		if sep {
			if fi.IsDir() {
				if !fj.IsDir() {
					return true
				}
			} else {
				if fj.IsDir() {
					return false
				}
			}
		}
		r := fi.ModTime().Unix() < fj.ModTime().Unix()
		if asc {
			return r
		}
		return !r
	})
	d.current = files
	return files
}

func HasPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return !(fi.Mode()&os.ModeNamedPipe == 0)
}
