/*
Package sr contains schema registry related functions.

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
package sr

import (
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/descriptors.yaml
var descriptorsTransformation string

// NewDescriptorsCommand creates an initialized command object
func NewDescriptorsCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRListDescriptorsParams{}
	cmd := &cobra.Command{
		Use:                   "descriptors",
		Short:                 "Display descriptors (Schema Registry)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			p.SRDescriptorFormat = api.AcceptObjects
			helper.CheckErr(output.SetTransformationDesc(descriptorsTransformation))
			pager := helper.NewPager(api.SRListDescriptorsP, conf.Authentication, p.Request()).
				OF("results").PP("next").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlagsPaging(cmd)
	helper.AddPagingFlags(&p.PageParams, cmd.Flags())
	return cmd
}
