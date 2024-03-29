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

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewDescriptorCommand creates an initialized command object
func NewDescriptorCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRListDescriptorsParams{}
	cmd := &cobra.Command{
		Use:                   "descriptor",
		Short:                 "Display a descriptor",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MaximumNArgs(1),
		ValidArgs:             []string{"created"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())

			//helper.CheckErr(output.SetTransformationDesc(desc))
			helper.CheckErr(output.PrintResponse(api.SRGetDescriptor(context.Background(), conf.Authentication, p)))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
