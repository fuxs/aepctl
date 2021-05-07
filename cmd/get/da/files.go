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
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/files.yaml
var filesTransformation string

// NewDatasetsCommand creates an initialized command object
func NewFilesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	fc := &api.DAOptions{}
	cmd := &cobra.Command{
		Use:                   "files batchId",
		Short:                 "Display all datasets",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationDesc(filesTransformation))
			fc.ID = args[0]
			helper.CheckErr(output.Print(api.DAGetFilesP, conf.Authentication, fc.Params()))
		},
	}
	output.AddOutputFlags(cmd)
	addFlags(cmd, fc)
	return cmd
}

func addFlags(cmd *cobra.Command, da *api.DAOptions) {
	flags := cmd.Flags()
	flags.IntVarP(&da.Start, "start", "s", 0, "paging parameter")
	flags.IntVarP(&da.Limit, "limit", "l", 0, "limits the number of results")
}
