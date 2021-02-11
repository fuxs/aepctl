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
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type tokenTransformer struct{}

func (t *tokenTransformer) ToTable(i interface{}) (*util.Table, error) {
	table := util.NewTable([]string{"TOKEN"}, 2)
	table.HideHeader = true
	if token, ok := i.(*api.BearerToken); ok {
		table.Append(map[string]interface{}{"TOKEN": token.Token})
	}
	return table, nil
}

func (t *tokenTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	table := util.NewTable([]string{"TOKEN", "EXPIRES"}, 2)
	table.HideHeader = true
	if token, ok := i.(*api.BearerToken); ok {
		table.Append(map[string]interface{}{
			"TOKEN":   token.Token,
			"EXPIRES": token.LocalTime(),
		})
	}
	return table, nil
}

// NewTokenCommand creates an initialized command object
func NewTokenCommand(auth *helper.Authentication) *cobra.Command {
	output := helper.NewOutputConf(&tokenTransformer{})
	cmd := &cobra.Command{
		Use:  "token",
		Args: cobra.NoArgs,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(auth.Validate(), output.ValidateFlags())
			output.PrintResult(auth.Config.GetToken())
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
