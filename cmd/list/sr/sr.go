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
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed trans/short.yaml
var shortTransformation string

// NewStatsCommand creates an initialized command object
func newListCommand(conf *helper.Configuration, use, short, long, example string, f api.Func) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRListParams{}
	var (
		show bool
	)
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 short,
		Long:                  long,
		Example:               example,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if show {
				output.SetTransformation(helper.NewRefTransformer())
			} else {
				p.Short = true
				helper.CheckErr(output.SetTransformationDesc(shortTransformation))
			}
			pager := helper.NewPager(f, conf.Authentication, p.Params()).
				OF("results").PP("next").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.BoolVar(&show, "show", false, "Show resource definition")
	addTenantFlags(flags, &p.SRBaseParams)
	return cmd
}

func addTenantFlags(flags *pflag.FlagSet, p *api.SRBaseParams) {
	addFlags(flags, p)
	flags.BoolVar(&p.Global, "predefined", false, "return resources defined by Adobe")
}

func addFlags(flags *pflag.FlagSet, p *api.SRBaseParams) {
	flags.StringVar(&p.Property, "property", "", "Comma separated list of top-level object properties to be returned in the response")
	flags.StringVar(&p.OrderBy, "order", "", "Sort response by specified fields separated by \",\"")
	flags.StringVar(&p.Start, "start", "", "The start value of the first orderBy field")
	flags.UintVar(&p.Limit, "limit", 0, "Specify a limit for the number of results to be displayed")
}
