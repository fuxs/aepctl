/*
Package util util consists of general utility functions and structures.

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
package util

import (
	"net/url"
	"sort"
	"strings"
)

// Par takes a list of name value pairs and builds a url parameter list
func Par(pairs ...string) string {
	l := len(pairs)
	if l%2 != 0 {
		l--
	}
	if l == 0 {
		return ""
	}
	sb := strings.Builder{}
	sep := "?"
	for i := 0; i < l; {
		if pairs[i+1] != "" {
			sb.WriteString(sep)
			sb.WriteString(pairs[i])
			i++
			sb.WriteByte('=')
			sb.WriteString(url.QueryEscape(pairs[i]))
			i++
			sep = "&"
		} else {
			i = i + 2
		}

	}
	a := sb.String()
	return a
}

// GetParam returns the first parameter value with given name from the passed
// url
func GetParam(str, name string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.Query().Get(name), nil
}

type Params map[string][]string

func NewParams(pairs ...string) Params {
	result := make(Params, len(pairs)/2)
	result.Set(pairs...)
	return result
}

func (p Params) Add(name, value string) {
	p[name] = append(p[name], value)
}

func (p Params) Set(pairs ...string) {
	l := len(pairs)
	if l%2 != 0 {
		l--
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

func (p Params) Get(key string) string {
	s := p[key]
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func (p Params) Encode() string {
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
			buf.WriteByte(sep)
			sep = byte('&')
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

func (p Params) EncodeWithout(names ...string) string {
	if p == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(p))
	var found bool
	for k := range p {
		found = false
		for _, n := range names {
			if n == k {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	sep := byte('?')
	for _, k := range keys {
		vs := p[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			buf.WriteByte(sep)
			sep = byte('&')
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
