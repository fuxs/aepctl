package sr

import (
	"context"

	"github.com/fuxs/aepctl/api/sr"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

var yamlStatsRecCre = `
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

var yamlStats = `
iterator: filter
filter: [imsOrg, tenantId, counts]
columns:
  - name: ORG
    type: str
    path: [imsOrg]
  - name: TENANT
    type: str
    path: [tenantId]
  - name: "# SCHEMAS"
    type: num
    path: [counts, schemas]
  - name: "# MIXINS"
    type: num
    path: [counts, mixins]
  - name: "# DATATYPES"
    type: num
    path: [counts, datatypes]
  - name: "# CLASSES"
    type: num
    path: [counts, classes]
  - name: "# UNIONS"
    type: num
    path: [counts, unions]
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
		Args:                  cobra.MaximumNArgs(1),
		ValidArgs:             []string{"created"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			desc := yamlStats
			if len(args) == 1 {
				switch args[0] {
				case "created":
					desc = yamlStatsRecCre
				}
			}
			helper.CheckErr(output.SetTransformationDesc(desc))
			output.StreamResultRaw(sr.GetStatsRaw(context.Background(), conf.Authentication))
			return
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
