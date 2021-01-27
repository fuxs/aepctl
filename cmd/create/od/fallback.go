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

func prepareFallback(auth *helper.Authentication, fallback *od.Fallback) {
	oc := helper.NewPlacementsCache(auth.AC)
	ps := helper.NameToID(oc)
	cs := oc.Mapper([]string{"_instance", "@id"}, []string{"_instance", "xdm:channel"})
	ts := helper.NewNameToID(auth, od.TagSchema)
	for _, r := range fallback.Representations {
		for _, c := range r.Components {
			c.Type = helper.ContentSToL.GetL(c.Type)
		}
		r.Placement = ps.Get(r.Placement)
		if r.Channel == "" {
			r.Channel = cs[r.Placement]
		} else {
			r.Channel = helper.ChannelSToL.GetL(r.Channel)
		}
	}
	for i, t := range fallback.Tags {
		fallback.Tags[i] = ts.GetValue(t)
	}
}

// NewCreateFallbackCommand creates an initialized command object
func NewCreateFallbackCommand(auth *helper.Authentication) *cobra.Command {
	ac := auth.AC
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     "fallback",
		Aliases: []string{"fallbacks"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(ac.AutoFillContainer())
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					fallback := &od.Fallback{}
					if err := i.Load(fallback); err == nil {
						if fc.IsYAML() {
							prepareFallback(auth, fallback)
						}
						_, err = od.Create(context.Background(), auth.Config, ac.ContainerID, od.FallbackSchema, fallback)
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
	fc.AddMandatoryFileFlag(cmd)
	return cmd
}
