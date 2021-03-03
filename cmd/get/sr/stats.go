package sr

import (
	"context"

	"github.com/fuxs/aepctl/api/sr"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

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
			desc := pkger.Include("/trans/get/sr/stats.yaml")
			if len(args) == 1 {
				switch args[0] {
				case "created":
					desc = pkger.Include("/trans/get/sr/created.yaml")
				}
			}
			helper.CheckErr(output.SetTransformationFile(desc))
			output.StreamResultRaw(sr.GetStatsRaw(context.Background(), conf.Authentication))
			return
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
