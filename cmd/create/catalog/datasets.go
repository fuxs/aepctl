package catalog

import (
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewCreateActivityCommand creates an initialized command object
func NewCreateDatasetCommand(conf *helper.Configuration) *cobra.Command {
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:                   "dataset",
		Aliases:               []string{"datasets"},
		Short:                 "Create an offer decisioning activity",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
		},
	}
	fc.AddMandatoryFileFlag(cmd)
	return cmd
}
