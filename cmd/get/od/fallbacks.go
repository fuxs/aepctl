/*
Package od contains offer decisiong related functions.

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
package od

import (
	"io"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type fallbackTransformer struct{}

func (*fallbackTransformer) Header(wide bool) []string {
	return []string{"NAME", "STATUS", "LAST MODIFIED"}
}

func (*fallbackTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*fallbackTransformer) WriteRow(name string, q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	return w.Write(
		s.Str("xdm:name"),
		StatusMapper.Get(s.Str("xdm:status")),
		util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
	)
}

func (*fallbackTransformer) Iterator(io.ReadCloser) (util.JSONResponse, error) {
	return nil, nil
}

// NewFallbacksCommand creates an initialized command object
func NewFallbacksCommand(conf *helper.Configuration) *cobra.Command {
	ft := &fallbackTransformer{}
	return NewQueryCommand(
		conf,
		od.FallbackSchema,
		"fallbacks",
		ft)
}

// NewFallbackCommand creates an initialized command object
func NewFallbackCommand(conf *helper.Configuration) *cobra.Command {
	ft := &fallbackTransformer{}
	return NewGetCommand(
		conf,
		helper.NewFallbackIDCache(conf.AC),
		od.FallbackSchema,
		"fallback",
		ft)
}
