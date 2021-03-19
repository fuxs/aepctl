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
package download

import (
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

func NewCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download dataset files",
	}
	conf.AddAuthenticationFlags(cmd)
	cmd.AddCommand(NewFilesCommand(conf))
	cmd.AddCommand(NewBatchesCommand(conf))
	cmd.AddCommand(NewDatasetsCommand(conf))
	return cmd
}
