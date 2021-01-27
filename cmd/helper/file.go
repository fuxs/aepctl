/*
Package helper consists of helping functions.

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
package helper

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Decoder is either a YAML or JSON decoder
type Decoder interface {
	Decode(v interface{}) error
}

// FileConfig provides the path and format of a file
type FileConfig struct {
	Path   string
	Format string
}

// AddFileFlag adds the required file flags to the passed command
func (c *FileConfig) AddFileFlag(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&c.Path, "file", "f", "", "a file")
	flags.StringVarP(&c.Format, "input", "i", "yaml", "the input format (yaml|yaml-raw)")
}

// AddMandatoryFileFlag adds the required file flags to the passed command. --file becomes mandatory.
func (c *FileConfig) AddMandatoryFileFlag(cmd *cobra.Command) {
	c.AddFileFlag(cmd)
	if err := cmd.MarkPersistentFlagRequired("file"); err != nil {
		fatal("Error in AddMandatoryFileFlag", 1)
	}
}

// IsYAML checks if format is YAML
func (c *FileConfig) IsYAML() bool {
	return strings.ToLower(c.Format) == "yaml"
}

// IsSet returns true if the path is set
func (c *FileConfig) IsSet() bool {
	return len(c.Path) > 0
}

// Open opens the file and returns a *FileIterator
func (c *FileConfig) Open() (*FileIterator, error) {
	if !c.IsSet() {
		return nil, nil
	}
	file := os.Stdin
	result := &FileIterator{}
	if c.Path != "-" {
		var err error
		if file, err = os.Open(c.Path); err != nil {
			return nil, err
		}
	}
	reader := bufio.NewReader(file)
	if strings.ToLower(filepath.Ext(c.Path)) == "json" {
		result.Decoder = json.NewDecoder(reader)
	} else {
		dec := yaml.NewDecoder(reader)
		//dec.SetStrict(true)
		result.Decoder = dec
	}
	return result, nil
}

// FileIterator helps iterating multiple YAML documents in one file
type FileIterator struct {
	Decoder Decoder
}

// Load loads the next available document in the file
func (c *FileIterator) Load(obj interface{}) error {
	if c.Decoder != nil {
		err := c.Decoder.Decode(obj)
		if err != nil {
			return err
		}
	}
	return nil
}
