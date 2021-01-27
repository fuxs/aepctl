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
	"github.com/spf13/cobra"
)

// NewCreateTagCommand creates an initialized command object
func NewCreateTagCommand(auth *helper.Authentication) *cobra.Command {
	ac := auth.AC
	ts := auth.TS
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     "tags",
		Aliases: []string{"tag"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(ac.AutoFillContainer())
			for i, name := range args {
				tag := &od.Tag{Name: name}
				_, err := od.Create(context.Background(), auth.Config, ac.ContainerID, od.TagSchema, tag)
				helper.CheckErr(err)
				if i == 0 {
					ts.Invalidate()
				}
			}

			if fc.IsSet() {
				i, err := fc.Open()
				helper.CheckErr(err)
				for {
					tag := &od.Tag{}
					err = i.Load(tag)
					if err == nil {
						_, err = od.Create(context.Background(), auth.Config, ac.ContainerID, od.TagSchema, tag)
						helper.CheckErr(err)
					} else {
						helper.CheckErrEOF(err)
						ts.Invalidate()
						break
					}
				}
			}
		},
	}
	ac.AddContainerFlag(cmd)
	fc.AddFileFlag(cmd)
	return cmd
}
