/*
Package get contains get command related functions.

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
package get

import (
	"time"

	"github.com/fuxs/aepctl/api/token"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type tokenTransformer struct{}

func (*tokenTransformer) Header(wide bool) []string {
	if wide {
		return []string{"TOKEN", "EXPIRES IN"}
	}
	return []string{"TOKEN"}
}

func (*tokenTransformer) Preprocess(i util.JSONResponse) error {
	return nil
}

func (*tokenTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	if wide {
		expires := q.Int("expires_in")
		return w.Write(q.Str("access_token"), time.Duration(expires*int(time.Millisecond)).String())
	}
	return w.Write(q.Str("access_token"))
}

// NewTokenCommand creates an initialized command object
func NewTokenCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(&tokenTransformer{})
	cmd := &cobra.Command{
		Use:  "token",
		Args: cobra.NoArgs,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			output.StreamResult(token.GetRaw(conf.Authentication))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
