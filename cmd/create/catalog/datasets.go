package catalog

import (
	"context"
	"fmt"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
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

// NewCreateActivityCommand creates an initialized command object
func NewCreateUnionProfileDatasetCommand(conf *helper.Configuration) *cobra.Command {
	var (
		format string
	)
	cmd := &cobra.Command{
		Use:                   "upds",
		Aliases:               []string{"union_profile", "unionProfile"},
		Short:                 "Create an dataset with the union profile schema",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			q, err := api.NewQuery(api.CatalogCreateProfileUnionDataset(context.Background(), conf.Authentication, args[0], format))
			helper.CheckErr(err)
			q.Range(func(q *util.Query) {
				fmt.Println(q.String())
			})
		},
	}
	flags := cmd.LocalFlags()
	flags.StringVarP(&format, "format", "f", "csv", "format of the data set file")
	return cmd
}
