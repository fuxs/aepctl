/*
Package cmd is the root package for aepctl.

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
package cmd

import (
	"github.com/fuxs/aepctl/cmd/completion"
	"github.com/fuxs/aepctl/cmd/create"
	"github.com/fuxs/aepctl/cmd/delete"
	"github.com/fuxs/aepctl/cmd/get"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/cmd/update"
	"github.com/fuxs/aepctl/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	longDesc = util.LongDesc(`
	A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`)
)

// NewCommand return an initialized command
func NewCommand() *cobra.Command {
	cobra.OnInitialize(initialze)
	cmd := &cobra.Command{
		Use:                   "aepctl",
		Short:                 "A brief description of your application",
		Long:                  longDesc,
		DisableFlagsInUseLine: true,
	}
	gcfg := util.NewGlobalConfig("aepctl", cmd)
	cmd.PersistentPreRunE = gcfg.GetPreRunE()
	auth := helper.NewAuthentication()
	auth.AddAuthenticationFlags(cmd)
	cmd.AddCommand(create.NewCommand(auth))
	cmd.AddCommand(get.NewCommand(auth))
	cmd.AddCommand(delete.NewCommand(auth))
	cmd.AddCommand(update.NewCommand(auth))
	cmd.AddCommand(completion.NewCommand())
	return cmd
}

func initialze() {
	// prepare loading of config file
	viper.SetEnvPrefix("mib")
	viper.AutomaticEnv()
}
