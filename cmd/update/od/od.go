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
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

func prepareUpdate(ac *cache.AutoContainer, update *od.Update, schema string) {
	idc := cache.NewODNameToInstanceIDMem(ac, schema)
	for i, name := range update.IDs {
		update.IDs[i] = idc.Lookup(name)
	}
}

// NewUpdateCommand creates an initialized update command object
func NewUpdateCommand(conf *helper.Configuration, ac *cache.AutoContainer, use, schema string) *cobra.Command {
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{util.Plural(use)},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					update := &od.Update{}
					if err := i.Load(update); err == nil {
						if fc.IsYAML() {
							prepareUpdate(ac, update, schema)
						}
						for _, name := range update.IDs {
							for _, apply := range update.Apply {
								_, err = od.Patch(context.Background(), conf.Authentication, cid, name, schema, apply)
								helper.CheckErr(err)
							}
						}

					} else {
						helper.CheckErrEOF(err)
						break
					}
				}
			}
		},
	}
	helper.CheckErr(ac.AddContainerFlag(cmd))
	fc.AddMandatoryFileFlag(cmd)
	return cmd
}
