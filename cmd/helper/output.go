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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
	"github.com/markbates/pkger"
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
	Iterator(io.ReadCloser) (util.JSONResponse, error)
}

/*type Pageable interface {
	InitialCall(context.Context, *api.AuthenticationConfig) (*http.Response, error)
	NextCall(context.Context, *api.AuthenticationConfig, string) (*http.Response, error)
}*/

// OutputConf contains all options for the output
type OutputConf struct {
	Output    string
	Type      OutputType
	jsonPath  string
	transPath string
	tf        Transformer
	td        *util.TableDescriptor
	//PB        Pageable
}

// NewOutputConf creates an initialized OutputConf object
func NewOutputConf(tf Transformer) *OutputConf {
	return &OutputConf{
		tf: tf,
	}
}

// SetTableTransformation changes the Transformer object
func (o *OutputConf) SetTableTransformation(td *util.TableDescriptor) {
	o.td = td
}

// SetTransformation changes the Transformer object
func (o *OutputConf) SetTransformation(tf Transformer) {
	o.tf = tf
}

// SetTransformationDesc changes the Transformer object
func (o *OutputConf) SetTransformationDesc(yaml string) error {
	tf, err := util.NewTableDescriptor(yaml)
	if err != nil {
		return err
	}
	o.tf = tf
	return nil
}

// SetTransformationFile changes the Transformer object
func (o *OutputConf) SetTransformationFile(path string) error {
	var (
		f   io.ReadCloser
		err error
		d   []byte
	)
	if o.transPath == "" {
		f, err = pkger.Open(path)
	} else {
		f, err = os.Open(o.transPath)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	if d, err = io.ReadAll(f); err != nil {
		return err
	}
	o.td, err = util.NewTableDescriptor(string(d))
	o.tf = o.td
	return err
}

// AddOutputFlags extends the passed command with flags for output
func (o *OutputConf) AddOutputFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&o.Output, "output", "o", "", "Output format (json|jsonpath=''|yaml|text).")
}

// ValidateFlags checks the passed flags
func (o *OutputConf) ValidateFlags() error {
	switch o.Output {
	case "", "table":
		o.Type = TableOut
	case "json":
		o.Type = JSONOut
	case "raw":
		o.Type = RawOut
	case "wide":
		o.Type = Wide
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
	CheckErr(err)
	if res.StatusCode >= 300 {
		data, err := io.ReadAll(res.Body)
		CheckErrs(err, errors.New(string(data)))
	}
	var (
		i util.JSONResponse
	)
	if o.tf != nil {
		i, err = o.tf.Iterator(res.Body)

	} else {
		i = util.NewJSONIterator(res.Body)
	}
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
	}
	return nil
}

func (o *OutputConf) Print(paged api.Paged) error {
	switch o.Type {
	case JSONOut:
		i, err := paged.First()
		CheckErr(err)
		return i.PrintPretty()
	case Wide, TableOut:
		return o.PrintTable(paged)
	}
	return nil
}

/*func (o *OutputConf) PrintJSON(auth *api.AuthenticationConfig) error {
	ctx := context.Background()
	res, err := o.PB.InitialCall(ctx, auth)
	CheckErr(err)
	if res.StatusCode >= 300 {
		data, err := io.ReadAll(res.Body)
		CheckErrs(err, errors.New(string(data)))
	}
	i, err := o.tf.Iterator(res.Body)
	CheckErr(err)
	return i.PrintPretty()
}*/

func (o *OutputConf) PrintTable(paged api.Paged) error {
	w := util.NewTableWriter(os.Stdout)
	defer func() {
		w.Flush()
	}()
	wide := o.Type == Wide
	if err := w.Write(o.tf.Header(wide)...); err != nil {
		return err
	}
	return paged.Execute(o.td.Path, func(j util.JSONResponse) error {
		return j.Range(func(q *util.Query) error {
			return o.tf.WriteRow(q, w, wide)
		})
	})
}
