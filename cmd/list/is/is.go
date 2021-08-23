/*
Package is contains identity service related functions.

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
package is

import (
	"context"
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/namespaces.yaml
var nsTransformation string

func NewNamespacesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	var imsOrg string
	cmd := &cobra.Command{
		Use:                   "namespaces",
		Short:                 "Display namespaces (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(nsTransformation))
			if imsOrg == "" {
				helper.CheckErr(output.PrintResponse(api.ISListNamespaces(context.Background(), conf.Authentication)))
			} else {
				helper.CheckErr(output.PrintResponse(api.ISListNamespacesIMSOrg(context.Background(), conf.Authentication, imsOrg)))
			}
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVar(&imsOrg, "ims-org", "", "IMS organization ID")
	return cmd
}
