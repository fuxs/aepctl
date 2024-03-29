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
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/sandboxes.yaml
var sandboxesTransformation string

//go:embed trans/details.yaml
var detailsTransformation string

//go:embed trans/types.yaml
var typesTransformation string

// NewSandboxCommand creates an initialized command object
func NewSandboxCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:  "sandbox",
		Args: cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if err := conf.Update(cmd); err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			sandboxes := cache.NewSandboxCache(conf.Authentication, conf).Values()
			return sandboxes, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			switch len(args) {
			case 0:
				helper.CheckErr(output.SetTransformationDesc(sandboxesTransformation))
				helper.CheckErr(output.PrintResponse(api.SBListSandboxes(context.Background(), conf.Authentication)))
			case 1:
				helper.CheckErr(output.SetTransformationDesc(detailsTransformation))
				helper.CheckErr(output.PrintResponse(api.SBGetSandbox(
					context.Background(),
					conf.Authentication,
					api.SBGetSandboxParams(args[0]))))
			}
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}

// NewSandboxesCommand creates an initialized command object
func NewSandboxesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:       "sandboxes",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"all", "types"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			switch len(args) {
			case 0:
				helper.CheckErr(output.SetTransformationDesc(sandboxesTransformation))
				helper.CheckErr(output.PrintResponse(api.SBListSandboxes(context.Background(), conf.Authentication)))
			case 1:
				switch args[0] {
				case "all":
					helper.CheckErr(output.SetTransformationDesc(sandboxesTransformation))
					helper.CheckErr(output.PrintResponse(api.SBListAllSandboxes(context.Background(), conf.Authentication)))
				case "types":
					helper.CheckErr(output.SetTransformationDesc(typesTransformation))
					helper.CheckErr(output.PrintResponse(api.SBListSandboxTypes(context.Background(), conf.Authentication)))
				}
			}
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
