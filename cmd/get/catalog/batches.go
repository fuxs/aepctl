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
	"context"
	_ "embed"
	"strconv"
	"time"

	"github.com/fuxs/aepctl/api/catalog"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/batches.yaml
var batchesTransformation string

// NewBatchesCommand creates an initialized command object
func NewBatchesCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	bc := &batchesConf{}
	cmd := &cobra.Command{
		Use:                   "batches",
		Short:                 "Display all batches",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			options, err := bc.ToOptions()
			helper.CheckErrs(err, output.SetTransformationDesc(batchesTransformation))
			output.StreamResultRaw(catalog.GetBatches(context.Background(), conf.Authentication, options))
		},
	}
	output.AddOutputFlags(cmd)
	bc.AddQueryFlags(cmd)
	return cmd
}

type batchesConf struct {
	name          string
	dataset       string
	limit         int
	timeFormat    string
	createdAfter  string
	createdBefore string
	startAfter    string
	startBefore   string
	endAfter      string
	endBefore     string
	orderBy       string
}

func (b *batchesConf) AddQueryFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVarP(&b.limit, "limit", "l", 0, "limits the number of results")
	flags.StringVar(&b.timeFormat, "time-format", time.RFC3339, "format for date parsing, default is '2006-01-02T15:04:05Z07:00' (RFC3339)")
	flags.StringVar(&b.createdAfter, "created-after", "", "returns batches created after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.createdBefore, "created-before", "", "returns batches created before this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.orderBy, "order-by", "", "sort parameter and direction for sorting, e.g. asc:created")
	flags.StringVar(&b.name, "name", "", "filter on the name of the dataset")
	flags.StringVar(&b.dataset, "dataset", "", "dataset identifier")
	flags.StringVar(&b.startAfter, "start-after", "", "returns batches started after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.startBefore, "start-before", "", "returns batches started before this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.endAfter, "end-after", "", "returns batches ended after this timestamp, see parameter time-format for encoding")
	flags.StringVar(&b.endBefore, "end-before", "", "returns batches ended before this timestamp, see parameter time-format for encoding")
}

func (b *batchesConf) ToOptions() (*catalog.BatchesOptions, error) {
	result := &catalog.BatchesOptions{
		Name:    b.name,
		Dataset: b.dataset,
		OrderBy: b.orderBy,
	}
	if b.limit > 0 {
		result.Limit = strconv.FormatInt(int64(b.limit), 10)
	}
	if b.createdAfter != "" {
		t, err := time.Parse(b.timeFormat, b.createdAfter)
		if err != nil {
			return nil, err
		}
		result.CreatedAfter = strconv.FormatInt(int64(t.Unix())*1000, 10)
	}
	if b.createdBefore != "" {
		t, err := time.Parse(b.timeFormat, b.createdBefore)
		if err != nil {
			return nil, err
		}
		result.CreatedBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.endAfter != "" {
		t, err := time.Parse(b.timeFormat, b.endAfter)
		if err != nil {
			return nil, err
		}
		result.EndAfter = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.endBefore != "" {
		t, err := time.Parse(b.timeFormat, b.endBefore)
		if err != nil {
			return nil, err
		}
		result.EndBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.startAfter != "" {
		t, err := time.Parse(b.timeFormat, b.startAfter)
		if err != nil {
			return nil, err
		}
		result.StartAfter = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.startBefore != "" {
		t, err := time.Parse(b.timeFormat, b.startAfter)
		if err != nil {
			return nil, err
		}
		result.StartBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	return result, nil
}
