/*
Package cmd is the root package for aepctl.

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
package download

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewODCommand creates an initialized command object
func NewBatchesCommand(conf *helper.Configuration) *cobra.Command {
	dc := &helper.DownloadConfig{}
	cmd := &cobra.Command{
		Use:                   "batches",
		Short:                 "Display all datasets",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd))
			bid := args[0]
			q, err := api.NewQuery(api.DAGetFiles(context.Background(), conf.Authentication, bid, "", ""))
			helper.CheckErr(err)
			q.Path("data").Range(func(q *util.Query) {
				fid := q.Str("dataSetFileId")
				q, err := api.NewQuery(api.DAGetFile(context.Background(), conf.Authentication, fid, "", ""))
				helper.CheckErr(err)
				q.Path("data").Range(func(q *util.Query) {
					name := q.Str("name")
					res, err := api.DADownload(context.Background(), conf.Authentication, fid, name)
					helper.CheckErr(err)
					helper.CheckErr(dc.Save(res, name))
				})
			})
		},
	}
	//output.AddOutputFlags(cmd)
	//fc.AddQueryFlags(cmd)
	return cmd
}
