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

type placementTransformer struct {
}

func (t *placementTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "CHANNEL TYPE", "CONTENT TYPE", "LAST MODIFIED", "DESCRIPTION"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          s.Str("xdm:name"),
			"CHANNEL TYPE":  helper.ChannelLToS.Get(s.Str("xdm:channel")),
			"CONTENT TYPE":  helper.ContentLToS.Get(s.Str("xdm:componentType")),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
			"DESCRIPTION":   s.Str("xdm:description"),
		})
	})
	return table, nil
}

func (t *placementTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"ID", "NAME", "CHANNEL TYPE", "CONTENT TYPE", "LAST MODIFIED", "DESCRIPTION"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"ID":            s.Str("@id"),
			"NAME":          s.Str("xdm:name"),
			"CHANNEL TYPE":  helper.ChannelLToS.Get(s.Str("xdm:channel")),
			"CONTENT TYPE":  helper.ContentLToS.Get(s.Str("xdm:componentType")),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
			"DESCRIPTION":   s.Str("xdm:description"),
		})
	})
	return table, nil
}

// NewPlacementsCommand creates an initialized command object
func NewPlacementsCommand(auth *helper.Authentication) *cobra.Command {
	pt := &placementTransformer{}
	return NewQueryCommand(
		auth,
		od.PlacementSchema,
		"placements",
		pt)
}

// NewPlacementCommand creates an initialized command object
func NewPlacementCommand(auth *helper.Authentication) *cobra.Command {
	pt := &placementTransformer{}
	return NewGetCommand(
		auth,
		helper.NewPlacementIDCache(auth.AC),
		od.PlacementSchema,
		"placement",
		pt)
}
