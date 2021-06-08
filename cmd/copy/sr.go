/*
Package copy is the base for all copy commands.

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
package copy

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewSRCommand creates an initialized command object
func NewSRCommand(conf *helper.Configuration) *cobra.Command {
	var destination string
	out := &helper.StatusConf{}
	cmd := &cobra.Command{
		Use:                   "sr",
		Short:                 "copy schema registry resource",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			destAuth := *conf.Authentication
			destAuth.Sandbox = destination
			// TODO mal testen
			for _, resource := range args {
				helper.CheckErr(conf.Validate(cmd))
				src, err := api.HandleStatusCode(api.SRExport(context.Background(), conf.Authentication, resource))
				helper.CheckErr(err)
				defer src.Body.Close()
				helper.CheckErr(out.PrintResponse(api.SRImportStream(context.Background(), &destAuth, src.Body)))
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&destination, "destination", "", "the destination sandbox")
	helper.CheckErr(cobra.MarkFlagRequired(flags, "destination"))
	return cmd
}
