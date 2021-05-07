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
package catalog

import (
	_ "embed"
	"time"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/batches.yaml
var batchesTransformation string

// NewBatchesCommand creates an initialized command object
func NewBatchesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	bc := &api.BatchesOptions{}
	cmd := &cobra.Command{
		Use:                   "batches",
		Short:                 "Display all batches",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(batchesTransformation))
			p := helper.CheckErrParams(bc)
			helper.CheckErr(output.Print(api.CatalogGetBatchesP, conf.Authentication, p))
		},
	}
	output.AddOutputFlags(cmd)
	addFlags(bc, cmd)
	return cmd
}

func addFlags(b *api.BatchesOptions, cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVarP(&b.Limit, "limit", "l", 0, "limits the number of results")
	flags.StringVar(&b.TimeFormat, "time-format", time.RFC3339, "format for date parsing, default is '2006-01-02T15:04:05Z07:00' (RFC3339)")
	flags.StringVar(&b.CreatedAfter, "created-after", "", "returns batches created after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.CreatedBefore, "created-before", "", "returns batches created before this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.OrderBy, "order-by", "", "sort parameter and direction for sorting, e.g. asc:created")
	flags.StringVar(&b.Name, "name", "", "filter on the name of the dataset")
	flags.StringVar(&b.Dataset, "dataset", "", "dataset identifier")
	flags.StringVar(&b.StartAfter, "start-after", "", "returns batches started after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.StartBefore, "start-before", "", "returns batches started before this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.EndAfter, "end-after", "", "returns batches ended after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.EndBefore, "end-before", "", "returns batches ended before this timestamp, see parameter time-format for encoding")
}
