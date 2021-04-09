/*
Package sr contains schema registry related functions.

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
package sr

import (
	"context"
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/schemas_sum.yaml
var schemasSumTransformation string

//go:embed trans/schemas_full.yaml
var schemasFullTransformation string

// NewStatsCommand creates an initialized command object
func NewSchemasCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	p := &api.SRGetSchemasParam{}
	cmd := &cobra.Command{
		Use:                   "schemas [global|tenant]",
		Short:                 "Display schemas",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MaximumNArgs(1),
		ValidArgs:             []string{"global", "tenant"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if len(args) == 1 {
				p.ContainerID = args[0]
			}
			desc := schemasSumTransformation
			if p.Full {
				desc = schemasFullTransformation
			}
			helper.CheckErr(output.SetTransformationDesc(desc))
			output.StreamResultRaw(api.SRGetSchemasRaw(context.Background(), conf.Authentication, p))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVar(&p.Properties, "properties", "", "Comma separated list of top-level object properties to be returned in the response")
	flags.StringVar(&p.OrderBy, "order", "", "Sort response by specified fields separated by \",\"")
	flags.StringVar(&p.Start, "start", "", "The start value of the first orderBy field")
	flags.UintVar(&p.Limit, "limit", 0, "Specify a limit for the number of results to be displayed")
	flags.BoolVar(&p.Full, "full", false, "Returns full JSON for each resource")
	return cmd
}