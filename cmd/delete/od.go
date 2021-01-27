/*
Package delete is the base for all delete commands.

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
package delete

import (
	"github.com/fuxs/aepctl/cmd/delete/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewODCommand creates an initialized command object
func NewODCommand(auth *helper.Authentication) *cobra.Command {
	cmd := &cobra.Command{
		Use: "od",
	}
	cmd.AddCommand(od.NewDeleteActivitiesCommand(auth))
	cmd.AddCommand(od.NewDeleteTagsCommand(auth))
	cmd.AddCommand(od.NewDeletePlacementsCommand(auth))
	cmd.AddCommand(od.NewDeleteOffersCommand(auth))
	cmd.AddCommand(od.NewDeleteFallbacksCommand(auth))
	cmd.AddCommand(od.NewDeleteCollectionsCommand(auth))
	cmd.AddCommand(od.NewDeleteRulesCommand(auth))
	return cmd
}
