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
	"context"
	"io"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// Pager executes a REST function and handles automatically the subsequent calls
// to get paged responses
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

// NewPager creates an initialzed Pager object. It uses JSONFilter to process
// payload and paging data.
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

// Next returns true, if a REST call can be executed
func (p *Pager) Next() bool {
	return p.calls == 0 || p.nextToken != ""
}

// Call executes the REST function
func (p *Pager) Call() error {
	if !p.Next() {
		return io.EOF
	}
	if p.calls == 0 {
		// first call
		if p.Context == nil {
			p.Context = context.Background()
		}
	} else {
		// subsequent call
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

// Run executes all REST calls
func (p *Pager) Run() error {
	for p.Next() {
		if err := p.Call(); err != nil {
			return err
		}
	}
	return nil
}

// Add adds the handler for the payload with the passed path
func (p *Pager) Add(f func(util.JSONResponse) error, path ...string) {
	p.jf.Add(f, path...)
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
