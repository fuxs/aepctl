package catalog

import (
	"context"
	"strconv"
	"time"

	"github.com/fuxs/aepctl/api/catalog"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type batchesConf struct {
	limit        int
	timeFormat   string
	createdAfter string
}

type batchTransformer struct{}

func (*batchTransformer) Header(wide bool) []string {
	return []string{"ID", "STATUS", "CREATED", "STARTED", "COMPLETED"}
}

func (*batchTransformer) Preprocess(i util.JSONResponse) error {
	return i.EnterObject()
}

// UTime
func uTime(q *util.Query) string {
	v := q.Integer()
	if v == 0 {
		return "-"
	}
	return time.Unix(int64(v)/1000, 0).Local().Format(time.RFC822)
}

func (*batchTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	r := q.Get(1)
	return w.Write(
		q.Get(0).String(),
		r.Str("status"),
		uTime(r.Path("created")),
		uTime(r.Path("started")),
		uTime(r.Path("completed")),
	)
}

// NewBatchesCommand creates an initialized command object
func NewBatchesCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(&batchTransformer{})
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
			helper.CheckErr(err)
			output.StreamResult(catalog.GetBatches(context.Background(), conf.Authentication, options))

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
}

func (b *batchesConf) ToOptions() (*catalog.BatchesOptions, error) {
	result := &catalog.BatchesOptions{}
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
