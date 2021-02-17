/*
Package get contains get command related functions.

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
package get

import (
	"context"
	"errors"
	"fmt"

	"github.com/fuxs/aepctl/api/sandbox"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type sandboxTransformer struct {
}

func (t *sandboxTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_page", "count")
	table := util.NewTable([]string{"NAME", "TITLE", "TYPE", "LAST MODIFIED"}, capacity)
	query.Path("sandboxes").Range(func(q *util.Query) {
		table.Append(map[string]interface{}{
			"NAME":          q.Str("name"),
			"TITLE":         q.Str("title"),
			"TYPE":          q.Str("type"),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("lastModifiedDate")),
		})
	})
	return table, nil
}

func (t *sandboxTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	capacity := query.Int("_page", "count")
	table := util.NewTable([]string{"NAME", "TITLE", "TYPE", "LAST MODIFIED"}, capacity)
	query.Path("sandboxes").Range(func(q *util.Query) {
		table.Append(map[string]interface{}{
			"NAME":          q.Str("name"),
			"TITLE":         q.Str("title"),
			"TYPE":          q.Str("type"),
			"LAST MODIFIED": util.LocalTimeStr(q.Str("lastModifiedDate")),
		})
	})
	return table, nil
}

type sandboxTypeTransformer struct{}

func (t *sandboxTypeTransformer) ToTable(i interface{}) (*util.Table, error) {
	query := util.NewQuery(i)
	table := util.NewTable([]string{"NAME"}, 2)
	query.Path("sandboxTypes").Range(func(q *util.Query) {
		table.Append(map[string]interface{}{
			"NAME": q.String(),
		})
	})
	return table, nil
}

func (t *sandboxTypeTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	return t.ToTable(i)
}

type sandboxDetailTransformer struct{}

func (t *sandboxDetailTransformer) ToTable(i interface{}) (*util.Table, error) {
	q := util.NewQuery(i)
	table := util.NewTable([]string{"NAME", "TITLE", "TYPE", "STATE", "REGION"}, 2)

	table.Append(map[string]interface{}{
		"NAME":   q.Str("name"),
		"TITLE":  q.Str("title"),
		"TYPE":   q.Str("type"),
		"STATE":  q.Str("state"),
		"REGION": q.Str("region"),
	})
	return table, nil
}

func (t *sandboxDetailTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	return t.ToTable(i)
}

// NewSandboxCommand creates an initialized command object
func NewSandboxCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(&sandboxTransformer{})
	cmd := &cobra.Command{
		Use: "sandbox",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if err := conf.Update(cmd); err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			sandboxes, _ := helper.NewSandboxCache(conf).GetList()
			return sandboxes, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			switch len(args) {
			case 0:
				output.PrintResult(sandbox.List(context.Background(), conf.Authentication))
			case 1:
				output.SetTransformation(&sandboxDetailTransformer{})
				output.PrintResult(sandbox.Get(context.Background(), conf.Authentication, args[0]))
			default:
				return errors.New("Too many arguments")
			}
			return nil
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}

// NewSandboxesCommand creates an initialized command object
func NewSandboxesCommand(conf *helper.Configuration) *cobra.Command {
	output := helper.NewOutputConf(&sandboxTransformer{})
	cmd := &cobra.Command{
		Use:       "sandboxes",
		ValidArgs: []string{"all", "types"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			switch len(args) {
			case 0:
				output.PrintResult(sandbox.List(context.Background(), conf.Authentication))
			case 1:
				switch args[0] {
				case "all":
					output.PrintResult(sandbox.ListAll(context.Background(), conf.Authentication))
				case "types":
					output.SetTransformation(&sandboxTypeTransformer{})
					output.PrintResult(sandbox.ListTypes(context.Background(), conf.Authentication))
				default:
					return fmt.Errorf("Unknown argument %s", args[0])
				}
			default:
				return errors.New("Too many arguments")
			}
			return nil
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
