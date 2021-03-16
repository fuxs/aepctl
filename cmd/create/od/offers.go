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

func prepareOffer(ac *cache.AutoContainer, offer *od.Offer) {
	ps := cache.NewODNameToIDMem(ac, od.PlacementSchema)
	for _, r := range offer.Representations {
		for _, c := range r.Components {
			c.Type = helper.ContentSToL.GetL(c.Type)
		}
		if r.Channel != "" {
			r.Channel = helper.ChannelSToL.GetL(r.Channel)
		}
		r.Placement = ps.Lookup(r.Placement)
	}
	//rules
	rs := cache.NewODNameToIDMem(ac, od.RuleSchema)
	offer.Constraint.Rule = rs.Lookup(offer.Constraint.Rule)
	// tags
	ts := cache.NewODNameToIDMem(ac, od.TagSchema)
	for i, t := range offer.Tags {
		offer.Tags[i] = ts.Lookup(t)
	}
}

// NewCreateOfferCommand creates an initialized command object
func NewCreateOfferCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     "offer",
		Aliases: []string{"offers"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					offer := &od.Offer{}
					if err := i.Load(offer); err == nil {
						if fc.IsYAML() {
							prepareOffer(ac, offer)
						}
						_, err = od.Create(context.Background(), conf.Authentication, cid, od.OfferSchema, offer)
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
