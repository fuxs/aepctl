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
func NewODCommand(auth *helper.Authentication) *cobra.Command {
	cmd := &cobra.Command{
		Use: "od",
	}
	cmd.AddCommand(od.NewActivitiesCommand(auth))
	cmd.AddCommand(od.NewActivityCommand(auth))
	cmd.AddCommand(od.NewCollectionsCommand(auth))
	cmd.AddCommand(od.NewCollectionCommand(auth))
	cmd.AddCommand(od.NewFallbacksCommand(auth))
	cmd.AddCommand(od.NewFallbackCommand(auth))
	cmd.AddCommand(od.NewOffersCommand(auth))
	cmd.AddCommand(od.NewOfferCommand(auth))
	cmd.AddCommand(od.NewPlacementsCommand(auth))
	cmd.AddCommand(od.NewPlacementCommand(auth))
	cmd.AddCommand(od.NewRulesCommand(auth))
	cmd.AddCommand(od.NewRuleCommand(auth))
	cmd.AddCommand(od.NewTagsCommand(auth))
	cmd.AddCommand(od.NewTagCommand(auth))
	return cmd
}
