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
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

//go:embed trans/placements.yaml
var placementsTransformation string

// NewPlacementsCommand creates an initialized command object
func NewPlacementsCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	td, err := util.NewTableDescriptor(placementsTransformation)
	helper.CheckErr(err)
	//pt := &placementTransformer{}
	return NewQueryCommand(
		conf,
		ac,
		od.PlacementSchema,
		"placements",
		td)
}

// NewPlacementCommand creates an initialized command object
func NewPlacementCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	td, err := util.NewTableDescriptor(placementsTransformation)
	helper.CheckErr(err)
	//pt := &placementTransformer{}
	return NewGetCommand(
		conf,
		ac,
		od.PlacementSchema,
		"placement",
		td)
}
