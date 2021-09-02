/*
Package api is the base for all aep rest functions.

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
package api

import (
	"net/url"
	"sort"
	"strings"
)

// RequestHeader containts the Header values of an HTTP request, e.g. Accept: text/html
type RequestHeader map[string]string

// RequestQuery contains the Query values of an HTTP request, e.g. ?next=good
type RequestQuery map[string][]string

// NewRequestQuery accepts multiple query name value pairs
func NewRequestQuery(pairs ...string) RequestQuery {
	result := make(RequestQuery, len(pairs)/2)
	result.Set(pairs...)
	return result
}

func (p RequestQuery) Set(pairs ...string) {
	l := len(pairs)
	if l%2 != 0 {
		pairs = append(pairs, "")
		l++
	}
	if l == 0 {
		return
	}
	for i := 0; i < l; i = i + 2 {
		if pairs[i+1] != "" {
			p.Add(pairs[i], pairs[i+1])
		}
	}
}

func (p RequestQuery) Add(name, value string) {
	p[name] = append(p[name], value)
}

// Encode builds an URL query
func (p RequestQuery) Encode() string {
	if p == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, len(p))
	i := 0
	for k := range p {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	sep := byte('?')
	for _, k := range keys {
		vs := p[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if v != "" {
				buf.WriteByte(sep)
				sep = byte('&')
				buf.WriteString(keyEscaped)
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	return buf.String()
}

type RequestHost map[string]string

type Request struct {
	header RequestHeader
	query  RequestQuery
	body   []byte
	aux    map[string]string
}

func NewRequest(pairs ...string) *Request {
	return &Request{query: NewRequestQuery(pairs...)}
}

func NewRequestHeader(header RequestHeader, pairs ...string) *Request {
	return &Request{header: header, query: NewRequestQuery(pairs...)}
}

func NewRequestBody(body []byte, pairs ...string) *Request {
	return &Request{body: body, query: NewRequestQuery(pairs...)}
}

func (r *Request) EncodedQuery() string {
	return r.query.Encode()
}

func (r *Request) Header() RequestHeader {
	if r.header == nil {
		r.header = make(RequestHeader)
	}
	return r.header
}

func (r *Request) SetHeader(name, value string) {
	if r.header == nil {
		r.header = make(RequestHeader)
	}
	r.header[name] = value
}

func (r *Request) Accept(value string) {
	r.SetHeader("Accept", value)
}

func (r *Request) ContentType(value string) {
	r.SetHeader("Content-Type", value)
}

func (r *Request) AddQuery(name, value string) {
	if r.query == nil {
		r.query = make(RequestQuery)
	}
	r.query.Add(name, value)
}

func (r *Request) SetValue(name, value string) {
	if r.aux == nil {
		r.aux = make(map[string]string)
	}
	r.aux[name] = value
}

func (r *Request) GetValue(name string) string {
	return r.aux[name]
}

func (r *Request) GetValuePath(name string) string {
	return url.PathEscape(r.aux[name])
}

func (r *Request) GetValueQuery(name string) string {
	return url.QueryEscape(r.aux[name])
}

// GetValueV returns the value with the passed name. If the value doesn't exist
// it returns the second value.
func (r *Request) GetValueV(name, value string) string {
	if result, ok := r.aux[name]; ok {
		return result
	}
	return value
}
