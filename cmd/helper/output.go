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
	// RawOut is used for raw output
	RawOut
)

// Transformer objects will implement transformation logic for certain OutputTypes
type Transformer interface {
	Header(wide bool) []string
	Preprocess(util.JSONResponse) error
	WriteRow(*util.Query, *util.RowWriter, bool) error
}

// OutputConf contains all options for the output
type OutputConf struct {
	Output   string
	Type     OutputType
	jsonPath string
	tf       Transformer
}

// NewOutputConf creates an initialized OutputConf object
func NewOutputConf(tf Transformer) *OutputConf {
	return &OutputConf{
		tf: tf,
	}
}

// SetTransformation changes the Transformer object
func (o *OutputConf) SetTransformation(tf Transformer) {
	o.tf = tf
}

// AddOutputFlags extends the passed command with flags for output
func (o *OutputConf) AddOutputFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&o.Output, "output", "o", "", "Output format (json|jsonpath=''|yaml|text).")
}

// ValidateFlags checks the passed flags
func (o *OutputConf) ValidateFlags() error {
	switch o.Output {
	case "":
		o.Type = TableOut
	case "json":
		o.Type = JSONOut
	case "raw":
		o.Type = RawOut
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

// StreamResult prints the object in the desired output format
func (o *OutputConf) StreamResult(i util.JSONResponse, err error) {
	CheckErr(err)
	CheckErr(o.streamResult(i))
}

func (o *OutputConf) streamResult(i util.JSONResponse) error {
	switch o.Type {
	case RawOut:
		return i.PrintRaw()
	case JSONOut:
		return i.PrintPretty()
	case JSONPathOut:
		// unmarshall complete response
		v, err := i.Obj()
		if err != nil {
			return err
		}
		value, err := jsonpath.Get(o.jsonPath, v)
		if err != nil {
			return err
		}
		bout := bufio.NewWriter(os.Stdout)
		defer bout.Flush()
		enc := json.NewEncoder(bout)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	case Wide, TableOut:
		w := util.NewTableWriter(os.Stdout)
		defer func() {
			w.Flush()
			i.Close()
		}()
		if err := o.tf.Preprocess(i); err != nil {
			return err
		}
		wide := o.Type == Wide
		if err := w.Write(o.tf.Header(wide)...); err != nil {
			return err
		}
		for i.More() {
			obj, err := i.Next()
			if err != nil {
				return err
			}
			if err = o.tf.WriteRow(util.NewQuery(obj), w, wide); err != nil {
				return err
			}
		}
	}
	return nil
}
