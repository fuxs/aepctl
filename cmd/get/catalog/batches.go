package catalog

import (
	"context"
	"strconv"
	"time"

	"github.com/fuxs/aepctl/api/catalog"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

type batchesConf struct {
	name         string
	limit        int
	timeFormat   string
	createdAfter string
}

var yamlBatches = `
iterator: object
columns:
  - name: ID
    id: true
  - name: STATUS
    type: str
    path: [status]
  - name: CREATED
    type: num
    path: [created]
    format: utime	
  - name: STARTED
    type: num 
    path: [started]
    format: utime
  - name: COMPLETED
    type: num
    path: [completed]
    format: utime
`

// NewBatchesCommand creates an initialized command object
func NewBatchesCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	bc := &batchesConf{}
	cmd := &cobra.Command{
		Use:                   "batches",
		Short:                 "Display all batches",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			options, err := bc.ToOptions()
			helper.CheckErrs(err, output.SetTransformationDesc(yamlBatches))
			output.StreamResultRaw(catalog.GetBatches(context.Background(), conf.Authentication, options))

		},
	}
	output.AddOutputFlags(cmd)
	bc.AddQueryFlags(cmd)
	return cmd
}

func (b *batchesConf) AddQueryFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVarP(&b.limit, "limit", "l", 0, "limits the number of results")
	flags.StringVar(&b.timeFormat, "format", time.RFC3339, "format for date parsing, default is '2006-01-02T15:04:05Z07:00' (RFC3339)")
	flags.StringVar(&b.createdAfter, "created-after", "", "date in local format")
	flags.StringVar(&b.name, "name", "", "filter on the name of the dataset")
}

func (b *batchesConf) ToOptions() (*catalog.BatchesOptions, error) {
	result := &catalog.BatchesOptions{
		Name: b.name,
	}
	if b.limit > 0 {
		result.Limit = strconv.FormatInt(int64(b.limit), 10)
	}
	if b.createdAfter != "" {
		t, err := time.Parse(b.timeFormat, b.createdAfter)
		if err != nil {
			return nil, err
		}
		result.CreatedAfter = strconv.FormatInt(int64(t.Unix()), 10)
	}
	return result, nil
}
