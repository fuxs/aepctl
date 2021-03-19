package flow

import (
	"context"
	"net/http"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

// NewDatasetsCommand creates an initialized command object
func NewConnectionsCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(nil)
	cc := NewConnectionsConf()
	cmd := &cobra.Command{
		Use:                   "connections",
		Short:                 "Display all connectionss",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationFile(pkger.Include("/trans/get/flow/connections.yaml")))
			output.PB = cc
			helper.CheckErr(output.Print(conf.Authentication))
			/*r, err := api.NewRestreamable(api.FlowGetConnections(context.Background(), conf.Authentication, params))
			helper.CheckErr(err)
			q, err := util.NewQueryStream(r.Reader())
			href := q.Str("_links", "next", "href")

			helper.CheckErr(err)
			output.StreamResultRaw()*/
		},
	}
	output.AddOutputFlags(cmd)
	cc.AddFlags(cmd)
	return cmd
}

type connectionsConf struct {
	p *api.FlowGetConnectionsParams
}

func NewConnectionsConf() *connectionsConf {
	return &connectionsConf{p: &api.FlowGetConnectionsParams{}}
}

func (c *connectionsConf) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVar(&c.p.ContinuationToken, "token", "", "a token for fetching records for next page")
	f.BoolVar(&c.p.Count, "count", false, "boolean value specifying if the count of resources should be returned (true|false)")
	f.StringVar(&c.p.Property, "property", "", "comma separated list of top-level object properties to be returned")
	f.IntVar(&c.p.Limit, "limit", 0, "max number of objects to be returned")
	f.StringVar(&c.p.OrderBy, "order", "", "results will be sorted")
}

func (c *connectionsConf) InitialCall(ctx context.Context, auth *api.AuthenticationConfig) (*http.Response, error) {
	return api.FlowGetConnections(ctx, auth, c.p)
}

func (c *connectionsConf) NextCall(ctx context.Context, auth *api.AuthenticationConfig) (*http.Response, error) {
	return nil, nil
}
