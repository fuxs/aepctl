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
	"github.com/fuxs/aepctl/cmd/get/is"
	"github.com/fuxs/aepctl/cmd/get/sr"
	"github.com/fuxs/aepctl/cmd/helper"

	"github.com/spf13/cobra"
)

// NewCommand creates an initialized command object
func NewCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
	}
	conf.AddAuthenticationFlags(cmd)
	cmd.AddCommand(NewACCommand(conf))
	cmd.AddCommand(NewTokenCommand(conf))
	cmd.AddCommand(NewODCommand(conf))
	cmd.AddCommand(NewSandboxCommand(conf))
	cmd.AddCommand(NewSandboxesCommand(conf))
	cmd.AddCommand(NewCatalogCommand(conf))
	//
	// identity service commands
	cmd.AddCommand(is.NewNamespaceCommand(conf))
	cmd.AddCommand(is.NewClusterCommand(conf))
	cmd.AddCommand(is.NewClustersCommand(conf))
	cmd.AddCommand(is.NewHistoryCommand(conf))
	cmd.AddCommand(is.NewHistoriesCommand(conf))
	cmd.AddCommand(is.NewXIDCommand(conf))
	cmd.AddCommand(is.NewMappingCommand(conf))
	//
	// schema registry commands
	cmd.AddCommand(sr.NewBehaviorCommand(conf))
	cmd.AddCommand(sr.NewClassCommand(conf))
	cmd.AddCommand(sr.NewDataTypeCommand(conf))
	cmd.AddCommand(sr.NewDescriptorCommand(conf))
	cmd.AddCommand(sr.NewFieldGroupCommand(conf))
	cmd.AddCommand(sr.NewSampleCommand(conf))
	cmd.AddCommand(sr.NewStatsCommand(conf))
	cmd.AddCommand(sr.NewSchemaCommand(conf))
	cmd.AddCommand(sr.NewUnionCommand(conf))

	cmd.AddCommand(NewDataAccessCommand(conf))
	cmd.AddCommand(NewFlowCommand(conf))
	cmd.AddCommand(NewUPSCommand(conf))
	return cmd
}
