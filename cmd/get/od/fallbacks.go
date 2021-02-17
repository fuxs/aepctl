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

type fallbackTransformer struct {
}

func (t *fallbackTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":        s.Str("xdm:status"),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
}

func (t *fallbackTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_embedded", "count")
	table := util.NewTable([]string{"NAME", "STATUS", "LAST MODIFIED"}, capacity)

	query.Path("_embedded", "results").Range(func(q *util.Query) {
		s := q.Path("_instance")
		table.Append(map[string]interface{}{
			"NAME":          strings.Trim(s.Str("xdm:name"), " \t"),
			"STATUS":        s.Str("xdm:status"),
			"LAST MODIFIED": util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
		})
	})
	return table, nil
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
