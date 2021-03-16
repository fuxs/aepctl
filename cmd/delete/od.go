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
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/delete/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewODCommand creates an initialized command object
func NewODCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use: "od",
	}
	ac := cache.NewAutoContainer(conf.Authentication, conf)
	cmd.AddCommand(od.NewDeleteActivitiesCommand(conf, ac))
	cmd.AddCommand(od.NewDeleteTagsCommand(conf, ac))
	cmd.AddCommand(od.NewDeletePlacementsCommand(conf, ac))
	cmd.AddCommand(od.NewDeleteOffersCommand(conf, ac))
	cmd.AddCommand(od.NewDeleteFallbacksCommand(conf, ac))
	cmd.AddCommand(od.NewDeleteCollectionsCommand(conf, ac))
	cmd.AddCommand(od.NewDeleteRulesCommand(conf, ac))
	return cmd
}
