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
	"errors"
	"io"
	"net/http"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// Pager executes a REST function and handles automatically the subsequent calls
// to get paged responses
type Pager struct {
	Func         api.Func
	Auth         *api.AuthenticationConfig
	Values       util.Params
	Context      context.Context
	PageFilter   []string
	ObjectFilter []string
	PagePath     []string
	PageParams   []string
	nextParams   []string
	calls        int
	jf           *util.JSONFinder
}

// NewPager creates an initialzed Pager object. It uses JSONFilter to process
// payload and paging data.
func NewPager(f api.Func, auth *api.AuthenticationConfig, v util.Params) *Pager {
	result := &Pager{
		Func:         f,
		Auth:         auth,
		Values:       v,
		PageFilter:   []string{"_links"},
		ObjectFilter: []string{"items"},
		PagePath:     []string{"next", "href"},
		PageParams:   []string{"continuationToken"},
	}
	return result
}

// OF sets the object filter. The object is the payload of the JSON document
func (p *Pager) OF(path ...string) *Pager {
	p.ObjectFilter = path
	return p
}

// PF sets the page filter
func (p *Pager) PF(path ...string) *Pager {
	p.PageFilter = path
	return p
}

// PF sets the page path
func (p *Pager) PP(path ...string) *Pager {
	p.PagePath = path
	return p
}

// P sets the page parameters. These are the URL query parameters which are
// necessary for the following paging requests
func (p *Pager) P(params ...string) *Pager {
	p.PageParams = params
	return p
}

// Next returns true, if a REST call can be executed
func (p *Pager) Next() bool {
	return p.calls == 0 || len(p.nextParams) > 0
}

func (p *Pager) SingleCall() (*http.Response, error) {
	if p.calls > 0 {
		return nil, errors.New("not the first call")
	}
	if p.Context == nil {
		p.Context = context.Background()
	}
	return api.HandleStatusCode(p.Func(p.Context, p.Auth, p.Values))
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
		for i, n := range p.PageParams {
			p.Values[n] = []string{p.nextParams[i]}
		}
		p.nextParams = p.nextParams[:0]
	}
	res, err := api.HandleStatusCode(p.Func(p.Context, p.Auth, p.Values))
	if err != nil {
		return err
	}
	p.calls++
	i := util.NewJSONIterator(util.NewJSONCursor(res.Body))
	defer res.Body.Close()
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
func (p *Pager) SetObjectHandler(f func(util.JSONResponse) error) {
	p.jf.Add(f, p.ObjectFilter...)
}

func (p *Pager) Prepare() {
	if p.nextParams == nil {
		p.nextParams = make([]string, 0, len(p.PageParams))
	}
	jf := util.NewJSONFinder()
	jf.Add(func(j util.JSONResponse) error {
		q, err := j.Query()
		if err != nil {
			return err
		}
		qu := q.Path(p.PagePath...)
		if !qu.Nil() {
			url := qu.String()
			for _, name := range p.PageParams {
				value, err := util.GetParam(url, name)
				if err != nil {
					return err
				}
				p.nextParams = append(p.nextParams, value)
			}
		}
		return nil
	}, p.PageFilter...)
	p.jf = jf
}
