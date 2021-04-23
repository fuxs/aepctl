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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/fuxs/aepctl/api"
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
	// WideOut is used for wide tablized output
	WideOut
	// RawOut is used for raw output
	RawOut
	// PVOut is using two columns path and value
	PVOut
	// PVOut is using three columns name, value and path
	NVPOUT
)

// Transformer objects will implement transformation logic for certain OutputTypes
type Transformer interface {
	Header(wide bool) []string
	Preprocess(util.JSONResponse) error
	WriteRow(*util.Query, *util.RowWriter, bool) error
	Iterator(*util.JSONCursor) (util.JSONResponse, error)
}

// OutputConf contains all options for the output
type OutputConf struct {
	Output    string
	Type      OutputType
	jsonPath  string
	transPath string
	tf        Transformer
}

// SetTransformation changes the Transformer object
func (o *OutputConf) SetTransformation(tf Transformer) {
	o.tf = tf
}

// SetTransformationDesc changes the Transformer object
func (o *OutputConf) SetTransformationDesc(def string) error {
	yaml := def
	if o.transPath != "" {
		f, err := os.Open(o.transPath)
		if err != nil {
			return err
		}
		defer f.Close()
		var d []byte
		if d, err = ioutil.ReadAll(f); err != nil {
			return err
		}
		yaml = string(d)
	}
	td, err := util.NewTableDescriptor(yaml)
	//o.td = td
	o.tf = td
	return err
}

// AddOutputFlags extends the passed command with flags for output
func (o *OutputConf) AddOutputFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&o.Output, "output", "o", "", "Output format (json|jsonpath=''|nvp|pv|raw|table|wide)")
	if err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		return []string{"json", "jsonpath=", "nvp", "pv", "raw", "table", "wide"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fatal("Error in AddOutputFlags", 1)
	}
}

// ValidateFlags checks the passed flags
func (o *OutputConf) ValidateFlags() error {
	switch o.Output {
	case "", "table":
		o.Type = TableOut
	case "json":
		o.Type = JSONOut
	case "vp":
		o.Type = PVOut
	case "nvp":
		o.Type = NVPOUT
	case "raw":
		o.Type = RawOut
	case "wide":
		o.Type = WideOut
	default:
		switch {
		case strings.HasPrefix(o.Output, "table="):
			l := len(o.Output)
			tp := o.Output[6:l]
			o.transPath = util.RemoveQuotes(tp)
			o.Type = TableOut
		case strings.HasPrefix(o.Output, "wide="):
			l := len(o.Output)
			tp := o.Output[5:l]
			o.transPath = util.RemoveQuotes(tp)
			o.Type = TableOut
		case strings.HasPrefix(o.Output, "jsonpath="):
			l := len(o.Output)
			jp := o.Output[9:l]
			o.jsonPath = util.AddDollar(util.RemoveQuotes(jp))
			o.Type = JSONPathOut
		default:
			return fmt.Errorf("unknown output format %s", o.Output)
		}
	}
	return nil
}

func (o *OutputConf) StreamResultRaw(res *http.Response, err error) {
	res, err = api.HandleStatusCode(res, err)
	CheckErr(err)
	var (
		i util.JSONResponse
	)
	if o.tf == nil || o.Type == NVPOUT || o.Type == PVOut {
		o.tf = &util.NVPTransformer{}
	}
	i, err = o.tf.Iterator(util.NewJSONCursor(res.Body))
	CheckErr(err)
	CheckErr(o.streamResult(i))
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
		q, err := i.Next()
		if err != nil {
			return err
		}
		v := q.Interface()
		value, err := jsonpath.Get(o.jsonPath, v)
		if err != nil {
			return err
		}
		bout := bufio.NewWriter(os.Stdout)
		defer bout.Flush()
		enc := json.NewEncoder(bout)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	case NVPOUT, PVOut, WideOut, TableOut:
		w := util.NewTableWriter(os.Stdout)
		defer func() {
			w.Flush()
			i.Close()
		}()
		if err := o.streamTableHeader(w); err != nil {
			return err
		}
		if err := o.streamTableBody(i, w); err != nil {
			return err
		}
	}
	return nil
}

func (o *OutputConf) wide() bool {
	return o.Type == WideOut || o.Type == NVPOUT
}

func (o *OutputConf) streamTableHeader(w *util.RowWriter) error {
	return w.Write(o.tf.Header(o.wide())...)
}

func (o *OutputConf) streamTableBody(i util.JSONResponse, w *util.RowWriter) error {
	if err := o.tf.Preprocess(i); err != nil {
		return err
	}
	wide := o.wide()
	for i.More() {
		q, err := i.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err = o.tf.WriteRow(q, w, wide); err != nil {
			return err
		}
	}
	return nil
}

func (o *OutputConf) PrintPaged(pager *Pager) error {
	pager.Prepare()
	switch o.Type {
	case JSONOut:
		res, err := pager.SingleCall()
		if err != nil {
			return err
		}
		c := util.NewJSONCursor(res.Body)
		return c.PrintPretty()
	case WideOut, TableOut:
		return o.PrintTable(pager)
	}
	return nil
}

// PrintTable prints out multiple JSON responses into one table
func (o *OutputConf) PrintTable(pager *Pager) error {
	w := util.NewTableWriter(os.Stdout)
	defer w.Flush()
	// add JSON object handler
	pager.SetObjectHandler(func(j util.JSONResponse) error {
		// copy a reseted cursor
		c, err := j.Cursor().New()
		if err != nil {
			return err
		}
		// complete JSON document must be processed
		defer c.End()
		// create the new iterator for the copied cursor
		i, err := o.tf.Iterator(c)
		if err != nil {
			return err
		}
		// print table body
		return o.streamTableBody(i, w)
	})
	// print the header
	if err := o.streamTableHeader(w); err != nil {
		return err
	}
	// print the table body
	return pager.Run()
}
