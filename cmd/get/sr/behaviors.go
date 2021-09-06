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

// NewBehaviorCommand creates an initialized command object
func NewBehaviorCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRGetGlobalParams{}
	cmd := &cobra.Command{
		Use:                   "behavior",
		Short:                 "Display a behavior",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if p.Full {
				output.SetTransformation(helper.NewTreeTransformer("$"))
			} else {
				output.SetTransformation(helper.NewRefTransformer("$"))
			}
			p.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.SRGetBehavior(context.Background(), conf.Authentication, p)))
		},
	}
	output.AddOutputFlags(cmd)
	addAcceptVersionedFlags(cmd, &p.SRFormat)
	return cmd
}
