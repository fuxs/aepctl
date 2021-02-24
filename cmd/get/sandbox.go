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

type sandboxTransformer struct{}

func (*sandboxTransformer) Header(wide bool) []string {
	return []string{"NAME", "TITLE", "TYPE", "LAST MODIFIED"}
}

func (*sandboxTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("sandboxes"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*sandboxTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	return w.Write(
		q.Str("name"),
		q.Str("title"),
		q.Str("type"),
		util.LocalTimeStr(q.Str("lastModifiedDate")),
	)
}

type sandboxTypeTransformer struct{}

func (*sandboxTypeTransformer) Header(wide bool) []string {
	return []string{"NAME"}
}

func (*sandboxTypeTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("sandboxTypes"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*sandboxTypeTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	return w.Write(q.String())
}

type sandboxDetailTransformer struct{}

func (*sandboxDetailTransformer) Header(wide bool) []string {
	return []string{"NAME", "TITLE", "TYPE", "STATE", "REGION"}
}

func (*sandboxDetailTransformer) Preprocess(i util.JSONResponse) error {
	return nil
}

func (*sandboxDetailTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	return w.Write(
		q.Str("name"),
		q.Str("title"),
		q.Str("type"),
		q.Str("state"),
		q.Str("region"),
	)
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
				output.StreamResult(sandbox.ListRaw(context.Background(), conf.Authentication))
			case 1:
				output.SetTransformation(&sandboxDetailTransformer{})
				output.StreamResult(sandbox.GetRaw(context.Background(), conf.Authentication, args[0]))
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
				output.StreamResult(sandbox.ListRaw(context.Background(), conf.Authentication))
			case 1:
				switch args[0] {
				case "all":
					output.StreamResult(sandbox.ListAllRaw(context.Background(), conf.Authentication))
				case "types":
					output.SetTransformation(&sandboxTypeTransformer{})
					output.StreamResult(sandbox.ListTypesRaw(context.Background(), conf.Authentication))
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
