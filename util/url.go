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
