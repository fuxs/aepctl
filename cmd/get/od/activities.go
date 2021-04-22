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
	_ "embed"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/activities.yaml
var activitiesTransformation string

func getPlacementCache(ac *cache.AutoContainer) *cache.MapMemCache {
	t := cache.NewODTrans().K("_instance", "@id").V("_instance", "xdm:channel")
	c := cache.NewODCall(ac, od.PlacementSchema)
	return cache.NewMapMemCache(c, t)
}

// NewActivitiesCommand creates an initialized command object
func NewActivitiesCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewQueryCommand(
		conf,
		ac,
		od.ActivitySchema,
		"activities",
		activitiesTransformation, "placements", getPlacementCache(ac))
}

// NewActivityCommand creates an initialized command object
func NewActivityCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewGetCommand(
		conf,
		ac,
		od.ActivitySchema,
		"activity",
		activitiesTransformation, "placements", getPlacementCache(ac))
}
