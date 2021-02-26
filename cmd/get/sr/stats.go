package sr

import (
	"context"

	"github.com/fuxs/aepctl/api/sr"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

var yamlStats = `
path: [recentlyCreatedResources]
columns:
  - name: NAME
    type: str
    path: [title]
  - name: TYPE
    type: str
    path: [meta:resourceType]
  - name: CREATED
    type: str
    path: [meta:created]
`

// NewStatsCommand creates an initialized command object
func NewStatsCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	cmd := &cobra.Command{
		Use:                   "stats",
		Short:                 "Display all stats",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(yamlStats))
			output.StreamResult(sr.GetStats(context.Background(), conf.Authentication))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
