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
	"bytes"
	"fmt"
	"strings"

	"github.com/russross/blackfriday"
)

type textRenderer struct {
	listItemCount uint
	listLevel     uint
}

func (r *textRenderer) NormalText(out *bytes.Buffer, text []byte) {
	raw := string(text)
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && trimmed[0] != '_' {
			out.WriteString(" ")
		}
		out.WriteString(trimmed)
	}
}

func (r *textRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	r.listLevel++
	out.WriteString("\n")
	text()
	r.listLevel--
}

func (r *textRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if flags&blackfriday.LIST_ITEM_BEGINNING_OF_LIST != 0 {
		r.listItemCount = 1
	} else {
		r.listItemCount++
	}
	indent := strings.Repeat("  ", int(r.listLevel))
	var bullet string
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		bullet += fmt.Sprintf("%d.", r.listItemCount)
	} else {
		bullet += "*"
	}
	out.WriteString(indent + bullet + " ")
	out.Write(text)
	out.WriteByte('\n')
}

func (r *textRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	out.WriteByte('\n')
	text()
	out.WriteByte('\n')
}

// BlockCode renders a chunk of text that represents source code.
func (r *textRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	out.WriteByte('\n')
	lines := []string{}
	for _, line := range strings.Split(string(text), "\n") {
		indented := "  " + line
		lines = append(lines, indented)
	}
	out.WriteString(strings.Join(lines, "\n"))
}

func (r *textRenderer) GetFlags() int { return 0 }
func (r *textRenderer) HRule(out *bytes.Buffer) {
	out.WriteString("\n" + "----------" + "\n")
}
func (r *textRenderer) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n")
}

func (r *textRenderer) TitleBlock(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	text()
}

func (r *textRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) TableRow(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	out.Write(text)
}

func (r *textRenderer) TableCell(out *bytes.Buffer, text []byte, align int) {
	out.Write(text)
}

func (r *textRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	text()
}

func (r *textRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	out.Write(text)
}

func (r *textRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}

func (r *textRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) RawHtmlTag(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	out.Write(ref)
}

func (r *textRenderer) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (r *textRenderer) Smartypants(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (r *textRenderer) DocumentHeader(out *bytes.Buffer)                          {}
func (r *textRenderer) DocumentFooter(out *bytes.Buffer)                          {}
func (r *textRenderer) TocHeaderWithAnchor(text []byte, level int, anchor string) {}
func (r *textRenderer) TocHeader(text []byte, level int)                          {}
func (r *textRenderer) TocFinalize()                                              {}

func (r *textRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.Write(header)
	out.Write(body)
}

func (r *textRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.WriteRune('[')
	out.Write(link)
	out.WriteRune(']')
}

func (r *textRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	out.WriteRune('<')
	out.Write(link)
	out.WriteRune('>')
}
