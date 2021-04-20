/*
Package od contains offer decisiong related functions.

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
package flow

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/connections.yaml
var connectionsTransformation string

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
			helper.CheckErrs(output.SetTransformationDesc(connectionsTransformation))

			helper.CheckErr(output.Page(api.FlowGetConnections, conf.Authentication, cc.p.Params()))
			//output.PB = cc
			//fp := api.NewFlowPaged(context.Background(), conf.Authentication, cc.p)
			//helper.CheckErr(output.Print(fp))
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
	return api.FlowGetConnections(ctx, auth, c.p.Params())
}

func (c *connectionsConf) NextCall(ctx context.Context, auth *api.AuthenticationConfig, url string) (*http.Response, error) {
	return api.FlowGetNext(ctx, auth, url)
}
