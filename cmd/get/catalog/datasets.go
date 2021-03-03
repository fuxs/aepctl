/*
Package catalog contains catalog related functions.

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
package catalog

import (
	"context"

	"github.com/fuxs/aepctl/api/catalog"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

var yamlDatasets = `
iterator: object
columns:
  - name: NAME
    type: str
    path: [name]
  - name: CREATED
    type: num
    path: [created]
    format: utime
  - name: LAST BATCH STATUS
    type: str
    path: [lastBatchStatus]
  - name: LAST UPDATED
    type: num
    path: [updated]
    format: utime
`

// NewDatasetsCommand creates an initialized command object
func NewDatasetsCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	bc := &batchesConf{}
	cmd := &cobra.Command{
		Use:                   "datasets",
		Short:                 "Display all datasets",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			options, err := bc.ToOptions()
			helper.CheckErrs(err, output.SetTransformationDesc(yamlDatasets))
			output.StreamResultRaw(catalog.GetDatasets(context.Background(), conf.Authentication, options))
		},
	}
	output.AddOutputFlags(cmd)
	bc.AddQueryFlags(cmd)
	return cmd
}