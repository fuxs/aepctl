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
package od

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewGetCommand creates an initialized command object
func NewGetCommand(conf *helper.Configuration, ac *cache.AutoContainer, schema, use, t, n string, c *cache.MapMemCache) *cobra.Command {
	output := &helper.OutputConf{}
	cmd := &cobra.Command{
		Use:  use,
		Args: cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// TODO check Update function
			if err := conf.Update(cmd); err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			idc := cache.NewODNameToID(ac, use, schema, conf.Sandboxed())
			return util.Difference(idc.Keys(), args), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			idc := cache.NewODNameToID(ac, use, schema, conf.Sandboxed())
			idc.Delete()
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			td, err := util.NewTableDescriptor(t)
			helper.CheckErr(err)
			if c != nil {
				td.AddMapping(n, c.Mapper())
			}
			output.SetTransformation(td)
			cid, err := ac.Get()
			helper.CheckErr(err)
			gp := &api.ODGetParams{
				ContainerID: cid,
				Schema:      schema,
			}
			for _, name := range args {
				gp.ID = idc.Lookup(name)
				helper.CheckErr(output.PrintResponse(api.ODGet(context.Background(), conf.Authentication, gp)))
			}
		},
	}
	output.AddOutputFlags(cmd)
	helper.CheckErr(ac.AddContainerFlag(cmd))
	return cmd
}

// NewQueryCommand creates an initialized command object
func NewQueryCommand(conf *helper.Configuration, ac *cache.AutoContainer, schema, use, t, n string, c *cache.MapMemCache) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.ODQueryParames{Schema: schema}
	cmd := &cobra.Command{
		Use:  use,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			td, err := util.NewTableDescriptor(t)
			helper.CheckErr(err)
			if c != nil {
				td.AddMapping(n, c.Mapper())
			}
			output.SetTransformation(td)
			p.ContainerID, err = ac.Get()
			helper.CheckErr(err)
			pager := helper.NewPager(api.ODQueryP, conf.Authentication, p.Request()).
				OF("_embedded", "results").P("start", "orderby")
			helper.CheckErr(output.PrintPaged(pager))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.PersistentFlags()
	flags.StringVarP(&p.Query, "query", "q", "", "Query string to search for in selected fields")
	flags.StringVar(&p.QOP, "qop", "", "Applies AND or OR operator to values in q query string param.")
	flags.StringVarP(&p.Field, "field", "f", "", "List of fields to limit the search to")
	flags.StringVarP(&p.OrderBy, "order-by", "b", "", "Sort results by a specific property.")
	flags.IntVarP(&p.Limit, "limit", "l", 0, "limits the number of returned results per request")
	helper.CheckErr(ac.AddContainerFlag(cmd))
	return cmd
}
