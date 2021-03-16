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
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates an initialized command object
func NewDeleteCommand(conf *helper.Configuration, ac *cache.AutoContainer, name, schema string) *cobra.Command {
	use := util.Plural(name)
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{name},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if err := conf.Update(cmd); err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			idc := cache.NewODNameToInstanceID(ac, name, schema, conf.Sandboxed())
			return util.Difference(idc.Keys(), args), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			idc := cache.NewODNameToInstanceID(ac, name, schema, conf.Sandboxed())
			defer idc.Delete()
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			for _, name := range args {
				helper.CheckErr(od.Delete(context.Background(), conf.Authentication, cid, idc.Lookup(name)))
			}
		},
	}
	helper.CheckErr(ac.AddContainerFlag(cmd))
	return cmd
}
