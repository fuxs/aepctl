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

type activityTransformer struct {
	idStore *util.KVCache
}

func newActivityTransformer(auth *helper.Authentication) *activityTransformer {
	// get list of placements and store map[@id]channel
	store := helper.NewTemporaryCache(auth.AC, od.PlacementSchema, []string{"_instance", "@id"}, []string{"_instance", "xdm:channel"})
	return &activityTransformer{
		idStore: store,
	}
}

func (t *activityTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "START DATE", "END DATE", "CHANNEL TYPE", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		// load store and transform to short names
		t.idStore.MapValues(func(s string) string {
			return helper.ChannelLToS.Get(s)
		})
		table.Append(map[string]interface{}{
			"NAME":       strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":     StatusMapper.Get(s.Str("xdm:status")),
			"START DATE": util.LocalTimeStrCustom(s.Str("xdm:startDate"), shortDate),
			"END DATE":   util.LocalTimeStrCustom(s.Str("xdm:endDate"), shortDate),
			"CHANNEL TYPE": s.Path("xdm:criteria").Concat(",", func(q *util.Query) string {
				id := q.Path("xdm:placements").Get(0).String()
				return t.idStore.GetValue(id)
			}),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

func (t *activityTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "START DATE", "END DATE", "CHANNEL TYPE", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		// load store and transform to short names
		t.idStore.MapValues(func(s string) string {
			return helper.ChannelLToS.Get(s)
		})
		table.Append(map[string]interface{}{
			"NAME":       strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":     StatusMapper.Get(s.Str("xdm:status")),
			"START DATE": util.LocalTimeStrCustom(s.Str("xdm:startDate"), shortDate),
			"END DATE":   util.LocalTimeStrCustom(s.Str("xdm:endDate"), shortDate),
			"CHANNEL TYPE": s.Path("xdm:criteria").Concat(",", func(q *util.Query) string {
				id := q.Path("xdm:placements").Get(0).String()
				return t.idStore.GetValue(id)
			}),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

// NewActivitiesCommand creates an initialized command object
func NewActivitiesCommand(auth *helper.Authentication) *cobra.Command {
	at := newActivityTransformer(auth)
	return NewQueryCommand(
		auth,
		od.ActivitySchema,
		"activities",
		at)
}

// NewActivityCommand creates an initialized command object
func NewActivityCommand(auth *helper.Authentication) *cobra.Command {
	at := newActivityTransformer(auth)
	return NewGetCommand(
		auth,
		helper.NewActivityIDCache(auth.AC),
		od.ActivitySchema,
		"activity",
		at)
}
