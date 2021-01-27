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
package od

import (
	"context"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates an initialized command object
func NewDeleteCommand(auth *helper.Authentication, os *util.KVCache, name string) *cobra.Command {
	ac := auth.AC //helper.NewAutoContainer(auth)
	use := util.Plural(name)
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{name},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			valid, err := os.Keys()
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			return util.Difference(valid, args), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(ac.AutoFillContainer())
			for _, name := range args {
				helper.CheckErr(od.Delete(context.Background(), auth.Config, ac.ContainerID, os.GetValue(name)))
				os.Remove(name)
			}
		},
	}
	ac.AddContainerFlag(cmd)
	return cmd
}
