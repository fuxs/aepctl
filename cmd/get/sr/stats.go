package sr

import (
	"context"
	_ "embed"

	"github.com/fuxs/aepctl/api/sr"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/stats.yaml
var statsTransformation string

//go:embed trans/created.yaml
var createdTransformation string

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
			desc := statsTransformation
			if len(args) == 1 {
				switch args[0] {
				case "created":
					desc = createdTransformation
				}
			}
			helper.CheckErr(output.SetTransformationDesc(desc))
			output.StreamResultRaw(sr.GetStatsRaw(context.Background(), conf.Authentication))
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
