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
	"net/url"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/pflag"
)

func AddPagingFlags(params *api.PageParams, flags *pflag.FlagSet) {
	addPagingFlags(params, flags)
	flags.StringVar(&params.Start, "start", "", "start value of property specified by flag order")
}

func AddPagingFlagsToken(params *api.PageParams, flags *pflag.FlagSet) {
	addPagingFlags(params, flags)
	flags.StringVar(&params.Start, "token", "", "a token for fetching records for next page")
}

func addPagingFlags(params *api.PageParams, flags *pflag.FlagSet) {
	flags.StringVar(&params.Order, "order", "", "order the result either by property updated or created (default)")
	flags.IntVar(&params.Limit, "limit", -1, "limits the number of returned results per request")
	flags.StringVar(&params.Filter, "filter", "", "filter by property created, updated, state or id")
}

// Pager executes a REST function and handles automatically the subsequent calls
// to get paged responses
type Pager struct {
	Func         api.Func
	Auth         *api.AuthenticationConfig
	Requests     []*api.Request
	Context      context.Context
	PageFilter   []string
	ObjectFilter []string
	PagePath     []string
	PageParams   []string
	nextParams   []string
	requestNum   int
	calls        int
	jf           *util.JSONFinder
}

// NewPager creates an initialzed Pager object. It uses JSONFilter to process
// payload and paging data.
func NewPager(f api.Func, auth *api.AuthenticationConfig, requests ...*api.Request) *Pager {
	result := &Pager{
		Func:         f,
		Auth:         auth,
		Requests:     requests,
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

// PP sets the page path
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
	return p.calls == 0 || len(p.nextParams) > 0 || p.requestNum < (len(p.Requests)-1)
}

func (p *Pager) SingleCall() (*http.Response, error) {
	if p.calls > 0 {
		return nil, errors.New("not the first call")
	}
	if p.Context == nil {
		p.Context = context.Background()
	}
	if len(p.Requests) == 0 {
		return api.HandleStatusCode(p.Func(p.Context, p.Auth, nil))
	}
	return api.HandleStatusCode(p.Func(p.Context, p.Auth, p.Requests[0]))
}

// Call executes the REST function
func (p *Pager) Call() error {
	if !p.Next() {
		return io.EOF
	}
	var params *api.Request
	if p.calls == 0 {
		// first call
		if p.Context == nil {
			p.Context = context.Background()
		}
		if len(p.Requests) > 0 {
			params = p.Requests[0]
		}
	} else {
		// subsequent call
		if len(p.nextParams) > 0 {
			if len(p.Requests) > 0 {
				// use provided request
				params = p.Requests[p.requestNum].Clone()
			} else {
				params = api.NewRequest()
			}
			// add the paging parameters
			params.AddQueries(p.nextParams...)
			// clear slice and keep allocated memory
			p.nextParams = p.nextParams[:0]
		} else {
			p.requestNum++
			params = p.Requests[p.requestNum]
		}
	}
	res, err := api.HandleStatusCode(p.Func(p.Context, p.Auth, params))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	p.calls++
	i := util.NewJSONIterator(util.NewJSONCursor(res.Body))
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

// RunOnce executes only one REST call
func (p *Pager) RunOnce() error {
	if p.Next() {
		return p.Call()
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
	// this is the pager function
	jf.Add(func(j util.JSONResponse) error {
		q, err := j.Query()
		if err != nil {
			return err
		}
		// select the element with the paging information
		qu := q.Path(p.PagePath...)
		if !qu.Nil() {
			u, err := url.Parse(qu.String())
			if err != nil {
				return err
			}
			urlq := u.Query()
			// search for all
			for _, name := range p.PageParams {
				p.nextParams = append(p.nextParams, name, urlq.Get(name))
			}
		}
		return nil
	}, p.PageFilter...)
	p.jf = jf
}
