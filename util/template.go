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
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/russross/blackfriday"
)

// LongDesc normalizes the long description
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return (&normalizer{s}).heredoc().markdown().trim().string
}

// Example normalizes the example section
func Example(s string) string {
	if len(s) == 0 {
		return s
	}
	return (&normalizer{s}).heredoc().indent().string
}

type normalizer struct {
	string
}

func (n *normalizer) heredoc() *normalizer {
	n.string = heredoc.Doc(n.string)
	return n
}

func (n *normalizer) markdown() *normalizer {
	n.string = string(blackfriday.Markdown([]byte(n.string), &textRenderer{}, blackfriday.EXTENSION_NO_INTRA_EMPHASIS))
	return n
}

func (n *normalizer) trim() *normalizer {
	n.string = strings.TrimSpace(n.string)
	return n
}

func (n *normalizer) indent() *normalizer {
	var buffer strings.Builder
	for _, line := range strings.Split(n.string, "\n") {
		buffer.WriteString("  ")
		buffer.WriteString(line)
		buffer.WriteByte('\n')
	}
	n.string = buffer.String()
	return n
}
