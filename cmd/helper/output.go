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
	"context"
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

/*type Pageable interface {
	InitialCall(context.Context, *api.AuthenticationConfig) (*http.Response, error)
	NextCall(context.Context, *api.AuthenticationConfig, string) (*http.Response, error)
}*/

type Pager struct {
	Func      api.Func
	Auth      *api.AuthenticationConfig
	Values    util.Params
	Context   context.Context
	Filter    []string
	Path      []string
	Parameter string
	nextToken string
	calls     int
	jf        *util.JSONFinder
}

func NewPager(f api.Func, auth *api.AuthenticationConfig, v util.Params) *Pager {
	result := &Pager{
		Func:      f,
		Auth:      auth,
		Values:    v,
		Filter:    []string{"_links"},
		Path:      []string{"next", "href"},
		Parameter: "continuationToken",
	}
	result.initJF()
	return result
}

func (p *Pager) Next() bool {
	return p.calls == 0 || p.nextToken != ""
}

func (p *Pager) Call() error {
	if !p.Next() {
		return io.EOF
	}
	if p.calls == 0 {
		if p.Context == nil {
			p.Context = context.Background()
		}
	}
	if p.nextToken != "" {
		p.Values[p.Parameter] = []string{p.nextToken}
		p.nextToken = ""
	}
	res, err := api.HandleStatusCode(p.Func(p.Context, p.Auth, p.Values))
	if err != nil {
		return err
	}
	p.calls++
	i := util.NewJSONIterator(util.NewJSONCursor(res.Body))
	defer i.Close()
	p.jf.SetIterator(i)

	return p.jf.Run()
}

func (p *Pager) initJF() {
	jf := util.NewJSONFinder()
	jf.Add(func(j util.JSONResponse) error {
		q, err := j.Query()
		if err != nil {
			return err
		}
		url := q.Str(p.Path...)
		token, err := util.GetParam(url, p.Parameter)
		if err != nil {
			return err
		}
		p.nextToken = token
		return nil
	}, p.Filter...)
	p.jf = jf
}

func (p *Pager) Add(f func(util.JSONResponse) error, path ...string) {
	p.jf.Add(f, path...)
}

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
	o.td = td
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

func (o *OutputConf) Print(paged api.Paged) error {
	switch o.Type {
	case JSONOut:
		i, err := paged.First()
		CheckErr(err)
		return i.PrintPretty()
	case WideOut, TableOut:
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
	wide := o.Type == WideOut
	if err := w.Write(o.tf.Header(wide)...); err != nil {
		return err
	}

	return paged.Execute(o.td.Path, func(j util.JSONResponse) error {
		return j.Range(func(q *util.Query) error {
			return o.tf.WriteRow(q, w, wide)
		})
	})
}

func (o *OutputConf) Page(f api.Func, auth *api.AuthenticationConfig, v util.Params) error {
	pager := NewPager(f, auth, v)
	switch o.Type {
	case WideOut, TableOut:

		return o.PrintTableP(pager)
	}
	return nil
}

func (o *OutputConf) PrintTableP(pager *Pager) error {
	w := util.NewTableWriter(os.Stdout)
	defer w.Flush()

	pager.Add(func(j util.JSONResponse) error {
		c, err := j.Cursor().New()
		if err != nil {
			return err
		}
		defer c.End()
		i, err := o.tf.Iterator(c)
		if err != nil {
			return err
		}
		return o.streamTableBody(i, w)
	}, o.td.ValuePath...)

	if err := o.streamTableHeader(w); err != nil {
		return err
	}

	for pager.Next() {
		if err := pager.Call(); err != nil {
			return err
		}
	}
	return nil
}
