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
	"fmt"
	"strings"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

func prepareCollection(ac *cache.AutoContainer, collection *od.Collection) {
	schema := od.OfferSchema
	if collection.Filter == "" {
		collection.Filter = "offers"
	} else {
		filter := helper.FilterTypeSToL.GetL(collection.Filter)
		collection.Filter = filter
		if filter == "anyTags" || filter == "allTags" {
			schema = od.TagSchema
		}
	}

	store := cache.NewODNameToIDMem(ac, schema)
	for i, o := range collection.IDs {
		collection.IDs[i] = store.Lookup(o)
	}
}

// NewCreateCollectionCommand creates an initialized command object
func NewCreateCollectionCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	fc := &helper.FileConfig{}

	cmd := &cobra.Command{
		Use:     "collection",
		Aliases: []string{"collections"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var result []string
			l := len(args)
			if l == 1 {
				result = []string{"all", "any", "offers"}
			} else if l > 1 {
				filter := strings.ToLower(args[1])

				if filter == "all" || filter == "any" {
					ts := cache.NewODNameToInstanceID(ac, "tags", od.TagSchema, conf.Sandboxed())
					result = ts.Keys()
				} else {
					os := cache.NewODNameToInstanceID(ac, "offers", od.TagSchema, conf.Sandboxed())
					result = os.Keys()
				}
			}
			if result == nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			return result, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			helper.CheckErr(conf.Validate(cmd))
			cid, err := ac.Get()
			helper.CheckErr(err)
			l := len(args)
			if l == 1 || l == 2 {
				return fmt.Errorf("invalid number of arguments (0, 3 or more): %v", l)
			}
			if l > 2 {
				filter := args[1]
				if filter != "all" && filter != "any" && filter != "offers" {
					return fmt.Errorf("second argument must be either any, all or offers: %v", filter)
				}
				collection := &od.Collection{
					Name:   args[0],
					Filter: filter,
					IDs:    args[2:],
				}
				prepareCollection(ac, collection)
				_, err := od.Create(context.Background(), conf.Authentication, cid, od.CollectionSchema, collection)
				helper.CheckErr(err)
			}
			i, err := fc.Open()
			helper.CheckErr(err)
			if i != nil {
				for {
					collection := &od.Collection{}
					if err := i.Load(collection); err == nil {
						if fc.IsYAML() {
							prepareCollection(ac, collection)
						}
						_, err = od.Create(context.Background(), conf.Authentication, cid, od.CollectionSchema, collection)
						helper.CheckErr(err)
					} else {
						helper.CheckErrEOF(err)
						break
					}
				}
			}
			return nil
		},
	}
	helper.CheckErr(ac.AddContainerFlag(cmd))
	fc.AddFileFlag(cmd)
	return cmd
}
