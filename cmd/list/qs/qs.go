/*
Package qs contains query service related functions.

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
package qs

import (
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/queries.yaml
var qsTransformation string

func NewQueriesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	params := &api.QSListQueriesParams{}
	cmd := &cobra.Command{
		Use:                   "queries",
		Short:                 "List queries (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(qsTransformation))
			pager := helper.NewPager(api.QSListQueriesP, conf.Authentication, params.Params()).
				OF("queries").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	flags := cmd.Flags()
	flags.StringVar(&params.Order, "order", "", "order the result either by property updated or created (default)")
	flags.IntVar(&params.Limit, "limit", -1, "limits the number of returned results per request")
	flags.StringVar(&params.Start, "start", "", "start value of property specified by flag order")
	flags.StringVar(&params.Filter, "filter", "", "filter by property created, updated, state or id")
	flags.BoolVar(&params.ExcludeSoftDeleted, "exclude-deleted", true, "exclude queries that have been soft deleted")
	flags.BoolVar(&params.ExcludeHidden, "exclude-hidden", true, "exclude uninteresting queries")
	return cmd
}
