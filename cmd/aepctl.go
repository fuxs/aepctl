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
	"github.com/fuxs/aepctl/cmd/audit"
	"github.com/fuxs/aepctl/cmd/cancel"
	"github.com/fuxs/aepctl/cmd/completion"
	"github.com/fuxs/aepctl/cmd/configure"
	"github.com/fuxs/aepctl/cmd/create"
	"github.com/fuxs/aepctl/cmd/delete"
	"github.com/fuxs/aepctl/cmd/download"
	"github.com/fuxs/aepctl/cmd/export"
	"github.com/fuxs/aepctl/cmd/extern"
	"github.com/fuxs/aepctl/cmd/get"
	"github.com/fuxs/aepctl/cmd/helper"
	imp "github.com/fuxs/aepctl/cmd/import"
	"github.com/fuxs/aepctl/cmd/list"
	"github.com/fuxs/aepctl/cmd/patch"
	"github.com/fuxs/aepctl/cmd/update"
	"github.com/fuxs/aepctl/cmd/version"
	"github.com/fuxs/aepctl/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version will be filled with version info
	Version  = "unknown"
	longDesc = util.LongDesc(`
	The command line tool for AEP
	
	aepctl is a command line tool for the Adobe Experience Platform implementing a part of the REST API.`)
)

// NewCommand return an initialized command
func NewCommand() *cobra.Command {
	cobra.OnInitialize(initialze)
	cmd := &cobra.Command{
		Use:                   "aepctl",
		Short:                 "The command line tool for AEP",
		Long:                  longDesc,
		DisableFlagsInUseLine: true,
	}
	gcfg := util.NewRootConfig("aepctl", Version, cmd)
	conf := helper.NewConfiguration(gcfg)
	cmd.AddCommand(cancel.NewCommand(conf))
	cmd.AddCommand(create.NewCommand(conf))
	cmd.AddCommand(get.NewCommand(conf))
	cmd.AddCommand(delete.NewCommand(conf))
	cmd.AddCommand(list.NewCommand(conf))
	cmd.AddCommand(update.NewCommand(conf))
	cmd.AddCommand(patch.NewCommand(conf))
	cmd.AddCommand(completion.NewCommand())
	cmd.AddCommand(completion.NewZSHCommand(gcfg))
	cmd.AddCommand(configure.NewConfigureCommand(gcfg))
	cmd.AddCommand(download.NewCommand(conf))
	cmd.AddCommand(export.NewSRCommand(conf))
	cmd.AddCommand(imp.NewSRCommand(conf))
	cmd.AddCommand(audit.NewCommand(conf))
	cmd.AddCommand(version.NewCommand(Version))
	cmd.AddCommand(extern.NewPSQLCommand(conf))
	return cmd
}

func initialze() {
	// prepare loading of config file
	viper.SetEnvPrefix("mib")
	viper.AutomaticEnv()
}
