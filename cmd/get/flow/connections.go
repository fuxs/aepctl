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
package flow

import (
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/connections.yaml
var connectionsTransformation string

// NewDatasetsCommand creates an initialized command object
func NewConnectionsCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.FlowGetConnectionsParams{}
	cmd := &cobra.Command{
		Use:                   "connections",
		Short:                 "Display all connectionss",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationDesc(connectionsTransformation))
			pager := helper.NewPager(api.FlowGetConnectionsP, conf.Authentication, p.Params())
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlags(cmd)
	f := cmd.Flags()
	f.StringVar(&p.ContinuationToken, "token", "", "a token for fetching records for next page")
	f.BoolVar(&p.Count, "count", false, "boolean value specifying if the count of resources should be returned (true|false)")
	f.StringVar(&p.Property, "property", "", "comma separated list of top-level object properties to be returned")
	f.IntVar(&p.Limit, "limit", 0, "limits the number of returned results per request")
	f.StringVar(&p.OrderBy, "order", "", "results will be sorted")
	return cmd
}
