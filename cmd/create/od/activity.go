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

func prepareActivity(ac *cache.AutoContainer, activity *od.Activity) {
	ps := cache.NewODNameToIDMem(ac, od.PlacementSchema)
	cs := cache.NewODNameToIDMem(ac, od.CollectionSchema)
	fs := cache.NewODNameToIDMem(ac, od.FallbackSchema)
	for _, c := range activity.Criteria {
		for i, p := range c.Placements {
			c.Placements[i] = ps.Lookup(p)
		}
		c.Selection.Filter = cs.Lookup(c.Selection.Filter)
	}
	activity.Fallback = fs.Lookup(activity.Fallback)
}

var (
	activityLong = util.LongDesc(`
	Create an offer decisioning activity

	This command requires a YAML or JSON file with the following structure:
	
	name: <name of the activity>
	startDate: <start date>
	  endDate: <end date>
	  status: <draft or approved>
	  criteria:
	    - selection:
	        filter: <name or @id of a collection>
	      placements:
	        - <name or @id of a placement>
	  fallback: <name or @id of a fallback>
	
	`)

	activityExample = util.Example(`
	# Create activity from YAML file
	aepctl create od activity --file examples/activity.yaml

	# Create activity from heredoc
	aepctl create od activity --file - << EOF
	  name: aepctl - example activity
	  startDate: "2020-10-01T16:00:00Z"
	  endDate: "2020-10-01T16:00:00Z"
	  status: draft
	  criteria:
	    - selection:
	        filter: aepctl - example collection # use collection name or @id
	      placements:
	        - Web - Image                       # use placement name or @id
	  fallback: aeptctl - example fallback      # use fallback name or @id
	EOF
	`)
)

// NewCreateActivityCommand creates an initialized command object
func NewCreateActivityCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:                   "activity",
		Aliases:               []string{"activities"},
		Short:                 "Create an offer decisioning activity",
		Long:                  activityLong,
		Example:               activityExample,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					activity := &od.Activity{}
					if err := i.Load(activity); err == nil {
						if fc.IsYAML() {
							prepareActivity(ac, activity)
						}
						_, err = od.Create(context.Background(), conf.Authentication, cid, od.ActivitySchema, activity)
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
