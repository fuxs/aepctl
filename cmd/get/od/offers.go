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
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type offerTransformer struct{}

func (*offerTransformer) Header(wide bool) []string {
	return []string{"NAME", "STATUS", "PRIORITY", "START DATE", "END DATE", "LAST MODIFIED"}
}

func (*offerTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*offerTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	d := s.Path("xdm:selectionConstraint")
	return w.Write(
		s.Str("xdm:name"),
		StatusMapper.Lookup(s.Str("xdm:status")),
		s.Str("xdm:rank", "xdm:priority"),
		util.LocalTimeStrCustom(d.Str("xdm:startDate"), shortDate),
		util.LocalTimeStrCustom(d.Str("xdm:endDate"), shortDate),
		util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
	)
}

func (*offerTransformer) Iterator(io.ReadCloser) (util.JSONResponse, error) {
	return nil, nil
}

// NewOffersCommand creates an initialized command object
func NewOffersCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	ot := &offerTransformer{}
	return NewQueryCommand(
		conf,
		ac,
		od.OfferSchema,
		"offers",
		ot)
}

// NewOfferCommand creates an initialized command object
func NewOfferCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	ot := &offerTransformer{}
	return NewGetCommand(
		conf,
		ac,
		od.OfferSchema,
		"offer",
		ot)
}
