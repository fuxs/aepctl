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
	"strconv"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

// NewDatasetsCommand creates an initialized command object
func NewFilesCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	fc := &filesConf{}
	cmd := &cobra.Command{
		Use:                   "files batchId",
		Short:                 "Display all datasets",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationFile(pkger.Include("/trans/get/da/files.yaml")))
			start, limit := fc.Strings()
			output.StreamResultRaw(api.DAGetFiles(context.Background(), conf.Authentication, args[0], start, limit))
		},
	}
	output.AddOutputFlags(cmd)
	fc.AddQueryFlags(cmd)
	return cmd
}

type filesConf struct {
	start int
	limit int
}

func (b *filesConf) AddQueryFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVarP(&b.start, "start", "s", 0, "paging parameter")
	flags.IntVarP(&b.limit, "limit", "l", 0, "limits the number of results")
}

func (b *filesConf) Strings() (start string, limit string) {
	if b.start > 0 {
		start = strconv.FormatInt(int64(b.start), 10)
	}
	if b.limit > 0 {
		limit = strconv.FormatInt(int64(b.limit), 10)
	}
	return start, limit
}
