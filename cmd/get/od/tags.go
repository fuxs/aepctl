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
	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type tagTransformer struct {
}

func (t *tagTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          s.Str("xdm:name"),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
		})
	})
	return table, nil
}

func (t *tagTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          s.Str("xdm:name"),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
		})
	})
	return table, nil
}

// NewTagsCommand creates an initialized command object
func NewTagsCommand(auth *helper.Configuration) *cobra.Command {
	tt := &tagTransformer{}
	return NewQueryCommand(
		auth,
		od.TagSchema,
		"tags",
		tt)
}

// NewTagCommand creates an initialized command object
func NewTagCommand(auth *helper.Configuration) *cobra.Command {
	tt := &tagTransformer{}
	return NewGetCommand(
		auth,
		helper.NewTagIDCache(auth.AC),
		od.TagSchema,
		"tag",
		tt)
}
