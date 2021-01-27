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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// OutputType is used for the encoding of different output formats
type OutputType int

const (
	// JSONOut is used for JSON
	JSONOut OutputType = iota
	// JSONPathOut is used for JSON path
	JSONPathOut
	// YAMLOut is used for YAML
	YAMLOut
	// TableOut is used for tablized output
	TableOut
	// Wide is used for wide tablized output
	Wide
)

// OutputConf contains all options for the output
type OutputConf struct {
	Output   string
	Type     OutputType
	jsonPath string
	tm       Transformer
}

// NewOutputConf creates an initialized OutputConf object
func NewOutputConf(tf Transformer) *OutputConf {
	return &OutputConf{
		tm: tf,
	}
}

// SetTransformation changes the Transformer object
func (o *OutputConf) SetTransformation(tf Transformer) {
	o.tm = tf
}

// Transformer objects will implement transformation logic for certain OutputTypes
//type Transformer func(interface{}) (interface{}, error)
type Transformer interface {
	ToTable(interface{}) (*util.Table, error)
	ToWideTable(interface{}) (*util.Table, error)
}

// TransformerMap is a map of transformers
type TransformerMap map[OutputType]Transformer

// AddOutputFlags extends the passed command with flags for output
func (o *OutputConf) AddOutputFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&o.Output, "output", "o", "", "Output format (json|jsonpath=''|yaml|text).")
}

// PrintResult prints the object in the desired output format
func (o *OutputConf) PrintResult(obj interface{}, err error) {
	CheckErr(err)
	if obj != nil {
		CheckErr(o.printResult(obj))
	}

}

// PrintResult prints the object in the desired output format
func (o *OutputConf) printResult(obj interface{}) error {
	if o.Type == JSONPathOut {
		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		v := interface{}(nil)
		if err = json.Unmarshal(data, &v); err != nil {
			return err
		}

		value, err := jsonpath.Get(o.jsonPath, v)
		if err != nil {
			return err
		}

		data, err = json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	if o.Type == JSONOut {
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	if o.Type == TableOut {
		t, err := o.tm.ToTable(obj)
		if err != nil {
			return err
		}
		t.Print(os.Stdout)
		//return util.PrintTabbed(os.Stdout, t)
	}
	if o.Type == Wide {
		t, err := o.tm.ToWideTable(obj)
		if err != nil {
			return err
		}
		t.Print(os.Stdout)
	}
	return nil
}

// ValidateFlags checks the passed flags
func (o *OutputConf) ValidateFlags() error {
	switch o.Output {
	case "":
		o.Type = TableOut
	case "json":
		o.Type = JSONOut
	case "wide":
		o.Type = Wide
	default:
		if strings.HasPrefix(o.Output, "jsonpath=") {
			l := len(o.Output)
			jp := o.Output[9:l]
			o.jsonPath = util.AddDollar(util.RemoveQuotes(jp))
			o.Type = JSONPathOut
		} else {
			// exit prog
			return fmt.Errorf("Unknown output format %s", o.Output)
		}
	}
	return nil
}
