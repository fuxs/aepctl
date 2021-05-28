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
	"github.com/fuxs/aepctl/util"
)

type RefTransformer struct {
	Path []string
}

func NewRefTransformer(path ...string) *RefTransformer {
	return &RefTransformer{Path: path}
}

func (t *RefTransformer) Header(wide bool) []string {
	return []string{"STRUCTURE", "TYPE", "XDM TYPE"}
}

func (t *RefTransformer) Preprocess(i util.JSONResponse) error {
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

func (t *RefTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	if err := w.Write(q.Str("title"), q.Str("type"), q.Str("meta:xdmType")); err != nil {
		return err
	}
	allOf := q.Path("allOf")
	references := make([]*util.Query, 0, 16)
	var properties, required *util.Query
	allOf.Range(func(q *util.Query) {
		ref := q.Path("$ref")
		if !ref.Nil() {
			references = append(references, q)
		} else {
			ref = q.Path("properties")
			if !ref.Nil() {
				properties = ref
			} else {
				ref = q.Path("required")
				if !ref.Nil() {
					required = ref
				}
			}
		}
	})

	definitions := q.Path("definitions")
	extends := q.Path("meta:extends")
	el := extends.Length()
	descriptors := q.Path("meta:descriptors")
	del := descriptors.Length()
	if len(references) > 0 {
		if !properties.Nil() || !required.Nil() || !definitions.Nil() || el > 0 || del > 0 {
			for _, q := range references {
				if err := w.Write("├─> "+q.Str("$ref"), q.Str("type"), q.Str("meta:xdmType")); err != nil {
					return err
				}
			}
		} else {
			l := allOf.Length() - 1
			prefix := "├─> "
			for i, q := range references {
				if i == l {
					prefix = "└─> "
				}
				if err := w.Write(prefix+q.Str("$ref"), q.Str("type"), q.Str("meta:xdmType")); err != nil {
					return err
				}
			}
		}
	}
	if !properties.Nil() {
		if !required.Nil() || !definitions.Nil() || el > 0 || del > 0 {
			if err := w.Write("├── properties", "", ""); err != nil {
				return err
			}
		} else {
			if err := w.Write("└── properties", "", ""); err != nil {
				return err
			}
		}
		if err := rangeChildren(properties, w, "│   "); err != nil {
			return err
		}
	}
	if !required.Nil() {
		var prefix, last string
		if !definitions.Nil() || el > 0 || del > 0 {
			if err := w.Write("├── required", "", ""); err != nil {
				return err
			}
			prefix = "│   ├── "
			last = "│   └── "
		} else {
			if err := w.Write("└── required", "", ""); err != nil {
				return err
			}
			prefix = "    ├── "
			last = "    └── "
		}
		l := required.Length() - 1
		err := required.RangeIE(func(i int, q *util.Query) error {
			if i == l {
				prefix = last
			}
			return w.Write(prefix+q.String(), "", "")
		})
		if err != nil {
			return err
		}
	}

	if !definitions.Nil() {
		prefix := "│   "
		if el > 0 || del > 0 {
			if err := w.Write("├── definitions", "", ""); err != nil {
				return err
			}
		} else {
			if err := w.Write("└── definitions", "", ""); err != nil {
				return err
			}
			prefix = "    "
		}
		if wide {
			if err := rangeChildrenWide(definitions, w, prefix); err != nil {
				return err
			}
		} else {
			if err := rangeChildren(definitions, w, prefix); err != nil {
				return err
			}
		}
	}

	if el > 0 {
		var prefix, last string
		if del > 0 {
			if err := w.Write("├── extends", "", ""); err != nil {
				return err
			}
			prefix = "│   ├── "
			last = "│   └── "
		} else {
			if err := w.Write("└── extends", "", ""); err != nil {
				return err
			}
			prefix = "    ├── "
			last = "    └── "
		}
		l := el - 1
		err := extends.RangeIE(func(i int, q *util.Query) error {
			if i == l {
				prefix = last
			}
			return w.Write(prefix+q.String(), "", "")
		})
		if err != nil {
			return err
		}
	}
	mapper := &util.Mapper{"xdm:alternateDisplayInfo": "Display",
		"xdm:descriptorIdentity":          "Identitiy",
		"descriptorOneToOne":              "Relationship",
		"xdm:descriptorReferenceIdentity": "Referenceable",
	}
	if del > 0 {
		if err := w.Write("└── descriptors", "", ""); err != nil {
			return err
		}
		prefix := "    ├── "
		l := del - 1
		return descriptors.RangeIE(func(i int, q *util.Query) error {
			if i == l {
				prefix = "    └── "
			}
			return w.Write(prefix+q.Str("@id"), mapper.Lookup(q.Str("@type")), "")
		})
	}
	return nil
}

func (t *RefTransformer) Iterator(c *util.JSONCursor) (util.JSONResponse, error) {
	return util.NewJSONIterator(c), nil
}
