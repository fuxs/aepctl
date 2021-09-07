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
var qsQueriesTransformation string

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
			helper.CheckErr(output.SetTransformationDesc(qsQueriesTransformation))
			pager := helper.NewPager(api.QSListQueriesP, conf.Authentication, params.Request()).
				OF("queries").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	flags := cmd.Flags()
	helper.AddPagingFlags(&params.PageParams, flags)
	flags.BoolVar(&params.ExcludeSoftDeleted, "exclude-deleted", true, "exclude queries that have been soft deleted")
	flags.BoolVar(&params.ExcludeHidden, "exclude-hidden", true, "exclude uninteresting queries")
	return cmd
}

//go:embed trans/schedules.yaml
var qsSchedulesTransformation string

func NewSchedulesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	params := &api.PageParams{}
	cmd := &cobra.Command{
		Use:                   "schedules",
		Short:                 "List schedules (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(qsSchedulesTransformation))
			pager := helper.NewPager(api.QSListSchedulesP, conf.Authentication, params.Request()).
				OF("schedules").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	helper.AddPagingFlags(params, cmd.Flags())
	return cmd
}

//go:embed trans/runs.yaml
var qsRunsTransformation string

func NewRunsCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	params := &api.PageParams{}
	cmd := &cobra.Command{
		Use:                   "runs scheduleId",
		Short:                 "List of runs for scheduled query (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(qsRunsTransformation))
			req := params.Request()
			req.SetValue("id", args[0])
			pager := helper.NewPager(api.QSListRunsP, conf.Authentication, req).
				OF("runsSchedules").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	helper.AddPagingFlags(params, cmd.Flags())
	return cmd
}

//go:embed trans/templates.yaml
var qsTemplatesTransformation string

func NewTemplatesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	params := &api.PageParams{}
	cmd := &cobra.Command{
		Use:                   "templates",
		Short:                 "List query templates (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(qsTemplatesTransformation))
			pager := helper.NewPager(api.QSListTemplatesP, conf.Authentication, params.Request()).
				OF("templates").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	helper.AddPagingFlags(params, cmd.Flags())
	return cmd
}
