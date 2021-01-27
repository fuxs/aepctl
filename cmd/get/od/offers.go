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
	"strings"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type offerTransformer struct {
}

func (t *offerTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "PRIORITY", "START DATE", "END DATE", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		d := s.Path("xdm:selectionConstraint")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":        StatusMapper.Get(s.Str("xdm:status")),
			"PRIORITY":      s.Str("xdm:rank", "xdm:priority"),
			"START DATE":    util.LocalTimeStrCustom(d.Str("xdm:startDate"), shortDate),
			"END DATE":      util.LocalTimeStrCustom(d.Str("xdm:endDate"), shortDate),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

func (t *offerTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "PRIORITY", "START DATE", "END DATE", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		d := s.Path("xdm:selectionConstraint")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":        StatusMapper.Get(s.Str("xdm:status")),
			"PRIORITY":      s.Str("xdm:rank", "xdm:priority"),
			"START DATE":    util.LocalTimeStrCustom(d.Str("xdm:startDate"), shortDate),
			"END DATE":      util.LocalTimeStrCustom(d.Str("xdm:endDate"), shortDate),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

// NewOffersCommand creates an initialized command object
func NewOffersCommand(auth *helper.Authentication) *cobra.Command {
	ot := &offerTransformer{}
	return NewQueryCommand(
		auth,
		od.OfferSchema,
		"offers",
		ot)
}

// NewOfferCommand creates an initialized command object
func NewOfferCommand(auth *helper.Authentication) *cobra.Command {
	ot := &offerTransformer{}
	return NewGetCommand(
		auth,
		helper.NewOfferIDCache(auth.AC),
		od.OfferSchema,
		"offer",
		ot)
}
