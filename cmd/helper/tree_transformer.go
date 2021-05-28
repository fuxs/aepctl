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
	"strings"

	"github.com/fuxs/aepctl/util"
)

type TreeTransformer struct {
	Path []string
}

func NewTreeTransformer(path ...string) *TreeTransformer {
	return &TreeTransformer{Path: path}
}

func (t *TreeTransformer) Header(wide bool) []string {
	if wide {
		return []string{"STRUCTURE", "PATH", "VALUE"}
	}
	return []string{"STRUCTURE", "TYPE", "XDM TYPE"}
}

func (t *TreeTransformer) Preprocess(i util.JSONResponse) error {
	if len(t.Path) == 1 && t.Path[0] == "$" {
		return nil
	}
	if len(t.Path) > 0 {
		if err := i.Path(t.Path...); err != nil {
			return err
		}
	}
	return i.Enter()
}

func treePrefix(name, prefix string, last bool) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	if last {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	sb.WriteString(name)
	return sb.String()
}

func arrayPrefix(prefix string) string {
	return prefix + "└── *"
}

func additionalPrefix(prefix string) string {
	return prefix + "└── +"
}

func rangeArray(q *util.Query, w *util.RowWriter, prefix string) error {
	if q.Nil() {
		return nil
	}
	t := q.Str("type")
	if err := w.Write(arrayPrefix(prefix), t, q.Str("meta:xdmType")); err != nil {
		return err
	}
	prefix = prefix + "    "
	ref := q.Path("$ref")
	if !ref.Nil() {
		return w.Write(prefix+"└─> "+ref.String(), "", "")
	}
	if t == "array" {
		return rangeArray(q.Path("items"), w, prefix)
	} else {
		if err := rangeChildren(q.Path("properties"), w, prefix); err != nil {
			return err
		}
		return rangeAdditional(q.Path("additionalProperties"), w, prefix)
	}
}

func rangeAdditional(q *util.Query, w *util.RowWriter, prefix string) error {
	if q.Nil() {
		return nil
	}
	t := q.Str("type")
	if err := w.Write(additionalPrefix(prefix), t, q.Str("meta:xdmType")); err != nil {
		return err
	}
	prefix = prefix + "    "
	if t == "array" {
		return rangeArray(q.Path("items"), w, prefix)
	} else {
		if err := rangeChildren(q.Path("properties"), w, prefix); err != nil {
			return err
		}
		return rangeAdditional(q.Path("additionalProperties"), w, prefix)
	}
}

func rangeChildren(q *util.Query, w *util.RowWriter, prefix string) error {
	lp := prefix + "│   "
	mp := prefix + "    "
	return q.RangeSortedAttributesRichE(func(name string, q *util.Query, i, size int) error {
		last := (i == size-1)
		t := q.Str("type")
		if err := w.Write(treePrefix(name, prefix, last), t, q.Str("meta:xdmType")); err != nil {
			return err
		}
		np := mp
		if !last {
			np = lp
		}
		if t == "array" {
			return rangeArray(q.Path("items"), w, np)
		} else {
			if err := rangeChildren(q.Path("properties"), w, np); err != nil {
				return err
			}
			return rangeAdditional(q.Path("additionalProperties"), w, np)
		}
	})
}

func printDetails(q *util.Query, w *util.RowWriter, prefix string, exclude ...string) error {
	return q.ResetPath().RangeValuesE(func(q *util.Query) error {
		return w.Write(prefix, q.JSONFullPath(), strings.TrimSpace(q.String()))
	}, exclude...)
}

func printChildren(q *util.Query, w *util.RowWriter, prefix string) error {
	qi := q.Path("items")
	if !qi.Nil() {
		if err := printDetails(q, w, prefix+"│   ", "items"); err != nil {
			return err
		}
		return rangeArrayWide(qi, w, prefix)
	} else {
		qp := q.Path("properties")
		if !qp.Nil() {
			if err := printDetails(q, w, prefix+"│   ", "properties"); err != nil {
				return err
			}
			return rangeChildrenWide(qp, w, prefix)
		}
		qap := q.Path("additionalProperties")
		if !qap.Nil() {
			if err := printDetails(q, w, prefix+"│   ", "additionalProperties"); err != nil {
				return err
			}
			return rangeAdditionalWide(qap, w, prefix)
		}
		// no children
		return printDetails(q, w, prefix, "")
	}
}

func rangeAdditionalWide(q *util.Query, w *util.RowWriter, prefix string) error {
	if q.Nil() {
		return nil
	}
	if err := w.Write(additionalPrefix(prefix), "", ""); err != nil {
		return err
	}
	return printChildren(q, w, prefix+"    ")
}

func rangeArrayWide(q *util.Query, w *util.RowWriter, prefix string) error {
	if q.Nil() {
		return nil
	}
	if err := w.Write(arrayPrefix(prefix), "", ""); err != nil {
		return err
	}
	return printChildren(q, w, prefix+"    ")
}

func rangeChildrenWide(q *util.Query, w *util.RowWriter, prefix string) error {
	if q.Nil() {
		return nil
	}
	mp := prefix + "│   "
	lp := prefix + "    "
	return q.RangeSortedAttributesRichE(func(name string, q *util.Query, i, size int) error {
		last := (i == size-1)
		if err := w.Write(treePrefix(name, prefix, last), "", ""); err != nil {
			return err
		}
		np := mp
		if last {
			np = lp
		}
		return printChildren(q, w, np)
	})
}

func (t *TreeTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	if wide {
		if err := w.Write(q.Str("title"), "", ""); err != nil {
			return err
		}
		if err := printDetails(q, w, "│   ", "properties"); err != nil {
			return err
		}
		return rangeChildrenWide(q.Path("properties"), w, "")
	} else {
		if err := w.Write(q.Str("title"), q.Str("type"), q.Str("meta:xdmType")); err != nil {
			return err
		}
		return rangeChildren(q.Path("properties"), w, "")
	}
}

func (t *TreeTransformer) Iterator(c *util.JSONCursor) (util.JSONResponse, error) {
	return util.NewJSONIterator(c), nil
}
