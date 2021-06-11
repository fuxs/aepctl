/*
Package imp contains import command related functions.

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
package imp

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewSRCommand creates an initialized command object
func NewSRCommand(conf *helper.Configuration) *cobra.Command {
	out := &helper.StatusConf{}
	cmd := &cobra.Command{
		Use:                   "import resource_id",
		Short:                 "import schema registry resource",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			if util.HasPipe() || len(args) == 0 {
				resource, err := ioutil.ReadAll(os.Stdin)
				helper.CheckErr(err)
				helper.CheckErr(out.PrintResponse(api.SRImport(context.Background(), conf.Authentication, resource)))
			}
			for _, file := range args {
				resource, err := os.ReadFile(file)
				helper.CheckErr(err)
				helper.CheckErr(out.PrintResponse(api.SRImport(context.Background(), conf.Authentication, resource)))
			}
		},
	}
	conf.AddAuthenticationFlags(cmd)
	return cmd
}
