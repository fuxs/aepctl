/*
Package od contains offer decisiong related functions.

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
package da

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewDatasetsCommand creates an initialized command object
func NewFileCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	fc := &api.DAOptions{}
	cmd := &cobra.Command{
		Use:                   "file fileId",
		Short:                 "Display all datasets",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			fc.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.DAGetFile(context.Background(), conf.Authentication, fc)))
		},
	}
	output.AddOutputFlags(cmd)
	addFlags(cmd, fc)
	return cmd
}
