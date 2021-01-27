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

func prepareUpdate(auth *helper.Authentication, update *od.Update, schema string) {
	store := helper.NewNameToInstanceID(auth, schema)
	for i, name := range update.IDs {
		update.IDs[i] = store.GetValue(name)
	}
}

// NewUpdateCommand creates an initialized update command object
func NewUpdateCommand(auth *helper.Authentication, use, schema string) *cobra.Command {
	ac := auth.AC
	fc := &helper.FileConfig{}
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{util.Plural(use)},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(ac.AutoFillContainer())
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					update := &od.Update{}
					if err := i.Load(update); err == nil {
						if fc.IsYAML() {
							prepareUpdate(auth, update, schema)
						}
						for _, name := range update.IDs {
							for _, apply := range update.Apply {
								_, err = od.Patch(context.Background(), auth.Config, ac.ContainerID, name, schema, apply)
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
	ac.AddContainerFlag(cmd)
	fc.AddMandatoryFileFlag(cmd)
	return cmd
}
