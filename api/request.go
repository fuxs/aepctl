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

// Request consists of HTTP header, body, query and auxiliary values
type Request struct {
	header map[string]string
	query  map[string][]string
	body   []byte
	aux    map[string]string
	auxm   map[string][]string
}

// NewRequest accepts string pairs of name and value for the HTTP query
func NewRequest(queries ...string) *Request {
	req := &Request{}
	req.AddQueries(queries...)
	return req
}

// NewRequestHeader accepts string pairs of name and value for the HTTP header
func NewRequestHeader(header ...string) *Request {
	req := &Request{}
	req.SetHeaders(header...)
	return req
}

// NewRequestBody accpets a body payload and name value pairs for the HTTP query
func NewRequestBody(body []byte, queries ...string) *Request {
	req := &Request{body: body}
	req.AddQueries(queries...)
	return req
}

// NewRequestValues accepts pairs of name and value for multiple purposes, e.g.
// in the path or host
func NewRequestValues(values ...string) *Request {
	req := &Request{}
	req.SetValues(values...)
	return req
}

// Clone creates a copy with the same body array
func (r *Request) Clone() *Request {
	header := make(map[string]string, len(r.header))
	for k, v := range r.header {
		header[k] = v
	}
	query := make(map[string][]string, len(r.query))
	for k, v := range r.query {
		query[k] = v
	}
	aux := make(map[string]string, len(r.aux))
	for k, v := range r.aux {
		aux[k] = v
	}
	auxm := make(map[string][]string, len(r.auxm))
	for k, v := range r.auxm {
		auxm[k] = v
	}
	return &Request{
		header: header,
		query:  query,
		body:   r.body,
		aux:    aux,
		auxm:   auxm,
	}
}

func (r *Request) EncodedQuery() string {
	if r.query == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, len(r.query))
	i := 0
	for k := range r.query {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	sep := byte('?')
	for _, k := range keys {
		vs := r.query[k]
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

func (r *Request) Header() map[string]string {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	return r.header
}

func (r *Request) SetHeader(name, value string) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[name] = value
}

func (r *Request) SetHeaderIf(name, value string) {
	if r.header == nil {
		r.header = make(map[string]string)
		r.header[name] = value
		return
	}
	if _, ok := r.header[name]; ok {
		return
	}
	r.header[name] = value
}

func (p *Request) SetHeaders(values ...string) {
	l := len(values)
	if l%2 != 0 {
		values = append(values, "")
		l++
	}
	if l == 0 {
		return
	}
	for i := 0; i < l; i = i + 2 {
		if values[i+1] != "" {
			p.SetHeader(values[i], values[i+1])
		}
	}
}

func (r *Request) Accept(value string) {
	r.SetHeader("Accept", value)
}

func (r *Request) ContentType(value string) {
	r.SetHeader("Content-Type", value)
}

func (r *Request) AddQuery(name, value string) {
	if r.query == nil {
		r.query = make(map[string][]string)
	}
	r.query[name] = append(r.query[name], value)
}

func (p *Request) AddQueries(values ...string) {
	l := len(values)
	if l%2 != 0 {
		values = append(values, "")
		l++
	}
	if l == 0 {
		return
	}
	for i := 0; i < l; i = i + 2 {
		if values[i+1] != "" {
			p.AddQuery(values[i], values[i+1])
		}
	}
}

func (r *Request) SetValue(name, value string) {
	if r.aux == nil {
		r.aux = make(map[string]string)
	}
	r.aux[name] = value
}

func (r *Request) GetValue(name string) string {
	if r.aux != nil {
		return r.aux[name]
	}
	return ""
}

func (r *Request) GetValuePath(name string) string {
	if r.aux != nil {
		return url.PathEscape(r.aux[name])
	}
	return ""
}

func (r *Request) GetValueQuery(name string) string {
	if r.aux != nil {
		return url.QueryEscape(r.aux[name])
	}
	return ""
}

// GetValueV returns the value with the passed name. If the value doesn't exist
// it returns the second value.
func (r *Request) GetValueV(name, value string) string {
	if r.aux != nil {
		if result, ok := r.aux[name]; ok {
			return result
		}
		return value
	}
	return value
}

func (p *Request) SetValues(values ...string) {
	l := len(values)
	if l%2 != 0 {
		values = append(values, "")
		l++
	}
	if l == 0 {
		return
	}
	for i := 0; i < l; i = i + 2 {
		if values[i+1] != "" {
			p.SetValue(values[i], values[i+1])
		}
	}
}

func (r *Request) SetArray(name string, value []string) {
	if r.auxm == nil {
		r.auxm = make(map[string][]string)
	}
	r.auxm[name] = value
}

func (r *Request) GetArray(name string) []string {
	if r.auxm != nil {
		return r.auxm[name]
	}
	return []string{}
}
