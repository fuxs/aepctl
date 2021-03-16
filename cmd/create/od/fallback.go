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
	"github.com/spf13/cobra"
)

func prepareFallback(ac *cache.AutoContainer, fallback *od.Fallback) {
	call := cache.NewODCall(ac, od.PlacementSchema)
	ps := cache.NewMapMemCache(call, cache.ODNameToID())
	cs := cache.NewMapMemCache(call, cache.NewODTrans().K("_instance", "@id").V("_instance", "xdm:channel"))

	ts := cache.NewODNameToIDMem(ac, od.TagSchema)
	for _, r := range fallback.Representations {
		for _, c := range r.Components {
			c.Type = helper.ContentSToL.GetL(c.Type)
		}
		r.Placement = ps.Lookup(r.Placement)
		if r.Channel == "" {
			r.Channel = cs.Lookup(r.Placement)
		} else {
			r.Channel = helper.ChannelSToL.GetL(r.Channel)
		}
	}
	for i, t := range fallback.Tags {
		fallback.Tags[i] = ts.Lookup(t)
	}
}

// NewCreateFallbackCommand creates an initialized command object
func NewCreateFallbackCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     "fallback",
		Aliases: []string{"fallbacks"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					fallback := &od.Fallback{}
					if err := i.Load(fallback); err == nil {
						if fc.IsYAML() {
							prepareFallback(ac, fallback)
						}
						_, err = od.Create(context.Background(), conf.Authentication, cid, od.FallbackSchema, fallback)
						helper.CheckErr(err)
					} else {
						helper.CheckErrEOF(err)
						break
					}
				}
			}
		},
	}
	helper.CheckErr(ac.AddContainerFlag(cmd))
	fc.AddMandatoryFileFlag(cmd)
	return cmd
}
