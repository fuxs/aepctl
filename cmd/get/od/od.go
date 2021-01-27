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

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// StatusMapper maps status values to a pretty representation
var StatusMapper = util.Mapper{
	"live":     "● Live",
	"approved": "● Approved",
	"draft":    "◯ Draft",
}

const (
	shortDate = "01/02/2006"
	longDate  = "01/02/2006, 03:04 PM"
)

// QueryConf stores the values of the query command
type QueryConf struct {
	Query   string
	QOP     string
	Field   string
	OrderBy string
	Limit   string
}

// AddQueryFlags adds all flags offered by a query command
func (q *QueryConf) AddQueryFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&q.Query, "query", "q", "", "Query string to search for in selected fields")
	flags.StringVar(&q.QOP, "qop", "", "Applies AND or OR operator to values in q query string param.")
	flags.StringVarP(&q.Field, "field", "f", "", "List of fields to limit the search to")
	flags.StringVarP(&q.OrderBy, "order-by", "b", "", "Sort results by a specific property.")
	flags.StringVarP(&q.Limit, "limit", "l", "", "Limit the number of decision rules returned.")
}

// NewGetCommand creates an initialized command object
func NewGetCommand(auth *helper.Authentication, os *util.KVCache, schema, use string, t helper.Transformer) *cobra.Command {
	output := helper.NewOutputConf(t)
	ac := auth.AC //helper.NewAutoContainer(auth)
	cmd := &cobra.Command{
		Use:  use,
		Args: cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			valid, err := os.Keys()
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			return util.Difference(valid, args), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(ac.AutoFillContainer())
			for _, name := range args {
				helper.CheckErrs(output.ValidateFlags(), ac.AutoFillContainer())
				output.PrintResult(od.Get(context.Background(), auth.Config, ac.ContainerID, schema, os.GetValue(name)))
			}
		},
	}
	output.AddOutputFlags(cmd)
	ac.AddContainerFlag(cmd)
	return cmd
}

// NewQueryCommand creates an initialized command object
func NewQueryCommand(auth *helper.Authentication, schema, use string, t helper.Transformer) *cobra.Command {
	output := helper.NewOutputConf(t)
	qc := &QueryConf{}
	ac := auth.AC
	cmd := &cobra.Command{
		Use:  use,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(output.ValidateFlags(), ac.AutoFillContainer())
			output.PrintResult(od.Query(context.Background(), auth.Config, ac.ContainerID, schema, qc.Query, qc.QOP, qc.Field, qc.OrderBy, qc.Limit))
		},
	}
	output.AddOutputFlags(cmd)
	qc.AddQueryFlags(cmd)
	ac.AddContainerFlag(cmd)
	return cmd
}
