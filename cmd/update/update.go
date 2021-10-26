/*
Package update contains update command related functions.

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
package update

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
		Use:                   "update",
		Short:                 "Update a resource",
		Long:                  longDesc,
		DisableFlagsInUseLine: true,
	}
	conf.AddAuthenticationFlags(cmd)
	cmd.AddCommand(NewODCommand(conf))
	cmd.AddCommand(NewClassCommand(conf))
	cmd.AddCommand(NewDataTypeCommand(conf))
	cmd.AddCommand(NewDescriptorCommand(conf))
	cmd.AddCommand(NewFieldGroupCommand(conf))
	cmd.AddCommand(NewSchemaCommand(conf))
	cmd.AddCommand(NewNamespaceCommand(conf))
	cmd.AddCommand(NewQueryTemplateCommand(conf))
	cmd.AddCommand(NewScheduleCommand(conf))
	return cmd
}

// NewUpdateCommand creates an initialized command object
func NewUpdateCommand(conf *helper.Configuration, f api.FuncPostID, use, short, long, example string, aliases ...string) *cobra.Command {
	var (
		id       string
		response bool
	)
	cmd := &cobra.Command{
		Use:                   use,
		Aliases:               aliases,
		Short:                 short,
		Long:                  long,
		Example:               example,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			ctx := context.Background()
			fr := util.MultiFileReader{Files: args}
			if response {
				helper.CheckErr(fr.ReadAll(func(data []byte) error {
					return api.PrintResponse(f(ctx, conf.Authentication, id, data))
				}))
			} else {
				helper.CheckErr(fr.ReadAll(func(data []byte) error {
					return api.DropResponse(f(ctx, conf.Authentication, id, data))
				}))
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&response, "response", false, "Print out response")
	flags.StringVar(&id, "id", "", "@id of the resource")
	helper.CheckErr(cmd.MarkFlagRequired("id"))
	return cmd
}
