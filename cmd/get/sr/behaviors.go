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

//go:embed trans/behaviors.yaml
var behaviorsTransformation string

//go:embed trans/behavior.yaml
var behaviorTransformation string

// NewStatsCommand creates an initialized command object
func NewBehaviorsCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := api.SRFormat{}
	cmd := &cobra.Command{
		Use:                   "behaviors",
		Short:                 "Display behaviors",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(behaviorsTransformation))
			pager := helper.NewPager(api.SRGetBehaviorsP, conf.Authentication, p.Params()).
				OF("results").PP("next").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.BoolVar(&p.Short, "short", false, "returns s list of ids only")
	flags.BoolVar(&p.Full, "full", false, "$ref attributes and allOf will be resolved")
	flags.BoolVar(&p.NoText, "notext", false, "no titles or descriptions")
	return cmd
}

// NewStatsCommand creates an initialized command object
func NewBehaviorCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := api.SRBehaviorParams{}
	cmd := &cobra.Command{
		Use:                   "behavior",
		Short:                 "Display behavior",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(behaviorTransformation))
			p.ID = args[0]
			/*pager := helper.NewPager(api.SRGetBehaviorP, conf.Authentication, p.Params()).
			OF("results").PP("next").P("start", "orderby")*/
			helper.CheckErr(output.Print(api.SRGetBehaviorP, conf.Authentication, p.Params()))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.BoolVar(&p.Short, "short", false, "returns s list of ids only")
	flags.BoolVar(&p.Full, "full", false, "$ref attributes and allOf will be resolved")
	flags.BoolVar(&p.NoText, "notext", false, "no titles or descriptions")
	flags.StringVar(&p.Version, "version", "1", "major version of resource")
	return cmd
}
