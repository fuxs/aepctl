/*
Package export contains export command related functions.

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
package export

import (
	"context"
	"io"
	"os"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewStatsCommand creates an initialized command object
func NewSRCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sr",
		Short:                 "export schema registry resource",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			res, err := api.HandleStatusCode(api.SRExport(context.Background(), conf.Authentication, args[0]))
			helper.CheckErr(err)
			defer res.Body.Close()
			_, err = io.Copy(os.Stdout, res.Body)
			helper.CheckErr(err)
		},
	}
	return cmd
}
