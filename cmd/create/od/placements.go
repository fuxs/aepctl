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
	"fmt"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewCreatePlacementCommand creates an initialized command object
func NewCreatePlacementCommand(conf *helper.Configuration) *cobra.Command {
	ac := conf.AC
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     "placement",
		Aliases: []string{"placements"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 1:
				if fc.IsYAML() {
					return helper.ContentSToL.Keys(), cobra.ShellCompDirectiveNoFileComp
				}
				return helper.ContentLToS.Keys(), cobra.ShellCompDirectiveNoFileComp
			case 2:
				if fc.IsYAML() {
					return helper.ChannelSToL.Keys(), cobra.ShellCompDirectiveNoFileComp
				}
				return helper.ChannelLToS.Keys(), cobra.ShellCompDirectiveNoFileComp
			default:
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			helper.CheckErr(ac.AutoFillContainer())
			l := len(args)
			if l == 1 || l == 2 || l > 4 {
				helper.CheckErr(fmt.Errorf("Invalid number of arguments (0, 3 or 4): %v", l))
			}
			if l > 2 {
				helper.CheckErr(ac.AutoFillContainer())
				placement := &od.Placement{
					Name:    args[0],
					Content: args[1],
					Channel: args[2],
				}
				if l > 3 {
					placement.Description = args[3]
				}
				if fc.IsYAML() {
					placement.Channel = helper.ChannelSToL.GetL(placement.Channel)
					placement.Content = helper.ContentSToL.GetL(placement.Content)
				}
				_, err := od.Create(context.Background(), conf.Authentication, ac.ContainerID, od.PlacementSchema, placement)
				helper.CheckErr(err)
			}
			//
			// load file with placements
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					placement := &od.Placement{}
					if err := i.Load(placement); err == nil {
						if fc.IsYAML() {
							placement.Channel = helper.ChannelSToL.GetL(placement.Channel)
							placement.Content = helper.ContentSToL.GetL(placement.Content)
						}
						_, err = od.Create(context.Background(), conf.Authentication, ac.ContainerID, od.PlacementSchema, placement)
						helper.CheckErr(err)
					} else {
						helper.CheckErrEOF(err)
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
