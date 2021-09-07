/*
Package create is the base for all create commands.

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
package create

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

var (
	longDesc = util.LongDesc(`
	A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`)
)

// NewCommand creates an initialized command object
func NewCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create",
		Short:                 "Create a resource",
		Long:                  longDesc,
		DisableFlagsInUseLine: true,
	}
	conf.AddAuthenticationFlags(cmd)
	cmd.AddCommand(NewODCommand(conf))
	cmd.AddCommand(NewCatalogCommand(conf))
	cmd.AddCommand(NewNamespaceCommand(conf))
	cmd.AddCommand(NewClassCommand(conf))
	cmd.AddCommand(NewDataTypeCommand(conf))
	cmd.AddCommand(NewDescriptorCommand(conf))
	cmd.AddCommand(NewFieldGroupCommand(conf))
	cmd.AddCommand(NewScheduleCommand(conf))
	cmd.AddCommand(NewSchemaCommand(conf))
	return cmd
}

// NewCreateCommand creates an initialized command object
func NewCreateCommand(conf *helper.Configuration, f api.FuncPost, use, short, long, example string, aliases ...string) *cobra.Command {
	var response, ignore bool
	cmd := &cobra.Command{
		Use:                   use,
		Aliases:               aliases,
		Short:                 short,
		Long:                  long,
		Example:               example,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			ctx := context.Background()
			fr := util.MultiFileReader{Files: args}
			if response {
				helper.CheckErr(fr.ReadAll(func(data []byte) error {
					err := api.PrintResponse(f(ctx, conf.Authentication, data))
					if ignore {
						helper.CheckErrInfo(err)
						return nil
					}
					return err
				}))
			} else {
				helper.CheckErr(fr.ReadAll(func(data []byte) error {
					err := api.DropResponse(f(ctx, conf.Authentication, data))
					if ignore {
						helper.CheckErrInfo(err)
						return nil
					}
					return err
				}))
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&response, "response", false, "Print out response")
	flags.BoolVar(&ignore, "ignore", false, "Ignore errors (for multiple arguments)")
	return cmd
}
