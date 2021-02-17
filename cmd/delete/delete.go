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
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

var (
	longDesc = util.LongDesc(`
	A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`)
)

// NewCommand creates an initialized command object
func NewCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete",
		Short:                 "Delete one or many resources",
		Long:                  longDesc,
		DisableFlagsInUseLine: true,
	}
	conf.AddAuthenticationFlags(cmd)
	cmd.AddCommand(NewODCommand(conf))
	return cmd
}
