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
)

//go:embed trans/short.yaml
var shortTransformation string

type listOption int

const (
	// JSONOut is used for JSON
	ListPredefined listOption = iota
	ListCustom
	ListSelect
)

// NewStatsCommand creates an initialized command object
func newListCommand(conf *helper.Configuration, use, short, long, example string, f api.Func, o listOption) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRListParams{}
	all := false
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
			// prepare two calls for flag --all
			params := make([]*api.Request, 1, 2)
			params[0] = p.Request()
			// show predefined and custom schemas?
			if all {
				// invert the Predefined flag
				p.SRBaseParams.Predefined = !p.SRBaseParams.Predefined
				params = append(params, p.Request())
			}
			pager := helper.NewPager(f, conf.Authentication, params...).
				OF("results").PP("next").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	flags := cmd.Flags()
	helper.AddPagingFlags(&p.PageParams, flags)
	flags.BoolVar(&show, "show", false, "Show resource definition")
	bp := &p.SRBaseParams
	switch o {
	case ListPredefined:
		bp.Predefined = true
	case ListCustom:
		bp.Predefined = false
	case ListSelect:
		flags.BoolVar(&bp.Predefined, "predefined", false, "return resources defined by Adobe")
		flags.BoolVar(&all, "all", false, "return resources defined by Adobe and custom resurces")
	}
	return cmd
}
