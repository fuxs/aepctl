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
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewStatsCommand creates an initialized command object
func newGetCommand(conf *helper.Configuration, use, short, long, example string, f api.Func) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.SRGetParams{}
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 short,
		Long:                  long,
		Example:               example,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if p.Full {
				output.SetTransformation(helper.NewTreeTransformer("$"))
			} else {
				output.SetTransformation(helper.NewRefTransformer("$"))
			}
			p.ID = args[0]
			helper.CheckErr(output.Print(f, conf.Authentication, p.Params()))
		},
	}
	output.AddOutputFlags(cmd)
	addAcceptVersionedFlags(cmd, &p.SRFormat)
	return cmd
}

func addAcceptVersionedFlags(cmd *cobra.Command, p *api.SRFormat) {
	flags := cmd.Flags()
	flags.BoolVar(&p.Short, "short", false, "returns s list of ids only")
	flags.BoolVar(&p.Full, "full", false, "$ref attributes and allOf will be resolved")
	flags.BoolVar(&p.NoText, "notext", false, "no titles or descriptions")
	flags.BoolVar(&p.Descriptors, "descriptors", false, "descriptors are included")
	flags.StringVar(&p.Version, "version", "1", "major version of resource")
}
