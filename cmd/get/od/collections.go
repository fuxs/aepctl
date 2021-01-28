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

type collectionTransformer struct {
}

func (t *collectionTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "# OFFERS", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"# OFFERS":      s.Len("xdm:ids"),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

func (t *collectionTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "# OFFERS", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"# OFFERS":      s.Len("xdm:ids"),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

// NewCollectionsCommand creates an initialized command object
func NewCollectionsCommand(auth *helper.Authentication) *cobra.Command {
	ct := &collectionTransformer{}
	return NewQueryCommand(
		auth,
		od.CollectionSchema,
		"collections",
		ct)
}

// NewCollectionCommand creates an initialized command object
func NewCollectionCommand(auth *helper.Authentication) *cobra.Command {
	ct := &collectionTransformer{}
	return NewGetCommand(
		auth,
		helper.NewCollectionIDCache(auth.AC),
		od.CollectionSchema,
		"collection",
		ct)
}