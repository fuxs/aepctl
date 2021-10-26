/*
Package trigger contains trigger command related functions.

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
package trigger

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewCancelRunCommand creates an initialized command object
func NewTriggerRunCommand(conf *helper.Configuration) *cobra.Command {
	var response, ignore bool
	cmd := &cobra.Command{
		Use:                   "run scheduleId",
		Short:                 "Triggers an immediate scheduled query run (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			ctx := context.Background()
			var err error
			for _, id := range args {
				if response {
					err = api.PrintResponse(api.QSTriggerRun(ctx, conf.Authentication, id))
				} else {
					err = api.DropResponse(api.QSTriggerRun(ctx, conf.Authentication, id))
				}
				if ignore {
					helper.CheckErrInfo(err)
				} else {
					helper.CheckErr(err)
				}
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&response, "response", false, "Print out response")
	flags.BoolVar(&ignore, "ignore", false, "Ignore errors (for multiple arguments)")
	return cmd
}
