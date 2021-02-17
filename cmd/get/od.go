/*
Package get contains get command related functions.

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
package get

import (
	"github.com/fuxs/aepctl/cmd/get/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewODCommand creates an initialized command object
func NewODCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use: "od",
	}
	cmd.AddCommand(od.NewActivitiesCommand(conf))
	cmd.AddCommand(od.NewActivityCommand(conf))
	cmd.AddCommand(od.NewCollectionsCommand(conf))
	cmd.AddCommand(od.NewCollectionCommand(conf))
	cmd.AddCommand(od.NewFallbacksCommand(conf))
	cmd.AddCommand(od.NewFallbackCommand(conf))
	cmd.AddCommand(od.NewOffersCommand(conf))
	cmd.AddCommand(od.NewOfferCommand(conf))
	cmd.AddCommand(od.NewPlacementsCommand(conf))
	cmd.AddCommand(od.NewPlacementCommand(conf))
	cmd.AddCommand(od.NewRulesCommand(conf))
	cmd.AddCommand(od.NewRuleCommand(conf))
	cmd.AddCommand(od.NewTagsCommand(conf))
	cmd.AddCommand(od.NewTagCommand(conf))
	return cmd
}
