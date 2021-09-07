/*
Package get contains get command related functions.

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
package get

import (
	"context"
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

func NewConnectionCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:                   "connection",
		Short:                 "Display connection parameters (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.PrintResponse(api.QSGetConnection(context.Background(), conf.Authentication)))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}

//go:embed trans/queries.yaml
var queryTransformation string

func NewQueryCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:                   "query",
		Short:                 "Display a query (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(queryTransformation))
			helper.CheckErr(output.PrintResponse(api.QSGetQuery(context.Background(), conf.Authentication, args[0])))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}

//go:embed trans/schedules.yaml
var scheduleTransformation string

func NewScheduleCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:                   "schedule",
		Short:                 "Display scheduled query (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(scheduleTransformation))
			helper.CheckErr(output.PrintResponse(api.QSGetSchedule(context.Background(), conf.Authentication, args[0])))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
