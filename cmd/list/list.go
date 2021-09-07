/*
Package list contains list command related functions.

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
package list

import (
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/cmd/list/is"
	"github.com/fuxs/aepctl/cmd/list/qs"
	"github.com/fuxs/aepctl/cmd/list/sr"
	"github.com/spf13/cobra"
)

// NewCommand creates an initialized command object
func NewCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List resources",
	}
	conf.AddAuthenticationFlags(cmd)
	//
	// quer service commands
	cmd.AddCommand(qs.NewQueriesCommand(conf))
	cmd.AddCommand(qs.NewRunsCommand(conf))
	cmd.AddCommand(qs.NewSchedulesCommand(conf))
	cmd.AddCommand(qs.NewTemplatesCommand(conf))
	//
	// identity service commands
	cmd.AddCommand(is.NewNamespacesCommand(conf))
	//
	// schema registry commands
	cmd.AddCommand(sr.NewBehaviorsCommand(conf))
	cmd.AddCommand(sr.NewClassesCommand(conf))
	cmd.AddCommand(sr.NewDataTypesCommand(conf))
	cmd.AddCommand(sr.NewDescriptorsCommand(conf))
	cmd.AddCommand(sr.NewFieldGroupsCommand(conf))
	cmd.AddCommand(sr.NewSchemasCommand(conf))
	cmd.AddCommand(sr.NewUnionsCommand(conf))
	return cmd
}
