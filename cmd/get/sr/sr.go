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
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewStatsCommand creates an initialized command object
func newQueryCommand(conf *helper.Configuration, use, short, long, example string, f api.Func) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRListParams{}
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 short,
		Long:                  long,
		Example:               example,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			desc := schemasSumTransformation
			if p.Full {
				desc = schemasFullTransformation
			}
			helper.CheckErr(output.SetTransformationDesc(desc))
			pager := helper.NewPager(f, conf.Authentication, p.Params()).
				OF("results").PP("next").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlags(cmd)
	addGetFlags(cmd, p)
	return cmd
}

// NewStatsCommand creates an initialized command object
func newGetCommand(conf *helper.Configuration, use, short, long, example, resource string, f api.Func) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRGetParams{}
	cmd := &cobra.Command{
		Use:                   "schema",
		Short:                 "Display schema",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			desc := schemasSumTransformation
			if p.Full {
				desc = schemasFullTransformation
			}
			p.ID = args[0]
			helper.CheckErr(output.SetTransformationDesc(desc))
			helper.CheckErr(output.Print(api.SRGetSchemaP, conf.Authentication, p.Params()))
		},
	}
	output.AddOutputFlags(cmd)
	addAcceptVersionedFlags(cmd, &p.SRFormat)
	return cmd
}

func addAcceptFlags(cmd *cobra.Command, p *api.SRFormat) {
	addAcceptStandard(cmd.Flags(), p)
}

func addAcceptVersionedFlags(cmd *cobra.Command, p *api.SRFormat) {
	flags := cmd.Flags()
	addAcceptStandard(flags, p)
	flags.StringVar(&p.Version, "version", "1", "major version of resource")
}

func addAcceptStandard(flags *pflag.FlagSet, p *api.SRFormat) {
	flags.BoolVar(&p.Short, "short", false, "returns s list of ids only")
	flags.BoolVar(&p.Full, "full", false, "$ref attributes and allOf will be resolved")
	flags.BoolVar(&p.NoText, "notext", false, "no titles or descriptions")
}

func addGetFlags(cmd *cobra.Command, p *api.SRListParams) {
	flags := cmd.Flags()
	addGetStandard(flags, &p.SRBaseParams)
	addAcceptStandard(flags, &p.SRFormat)
}

func addGetDescriptorsFlags(cmd *cobra.Command, p *api.SRListDescriptorsParams) {
	flags := cmd.Flags()
	addGetStandard(flags, &p.SRBaseParams)
	flags.BoolVar(&p.Short, "short", false, "returns s list of ids only")
	flags.BoolVar(&p.Full, "full", false, "$ref attributes and allOf will be resolved")
}

func addGetStandard(flags *pflag.FlagSet, p *api.SRBaseParams) {
	flags.StringVar(&p.Properties, "properties", "", "Comma separated list of top-level object properties to be returned in the response")
	flags.StringVar(&p.OrderBy, "order", "", "Sort response by specified fields separated by \",\"")
	flags.StringVar(&p.Start, "start", "", "The start value of the first orderBy field")
	flags.UintVar(&p.Limit, "limit", 0, "Specify a limit for the number of results to be displayed")
	flags.BoolVar(&p.Global, "global", false, "return core resources instead of custom resources")
}
