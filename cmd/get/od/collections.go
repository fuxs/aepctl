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
	_ "embed"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

/*type collectionTransformer struct{}

func (*collectionTransformer) Header(wide bool) []string {
	return []string{"NAME", "# OFFERS", "LAST MODIFIED"}
}

func (*collectionTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*collectionTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	return w.Write(
		s.Str("xdm:name"),
		strconv.Itoa(s.Len("xdm:ids")),
		util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
	)
}

func (*collectionTransformer) Iterator(c *util.JSONCursor) (util.JSONResponse, error) {
	return util.NewJSONIterator(c), nil
}*/

//go:embed trans/collections.yaml
var collectionsTransformation string

// NewCollectionsCommand creates an initialized command object
func NewCollectionsCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewQueryCommand(
		conf,
		ac,
		od.CollectionSchema,
		"collections",
		collectionsTransformation, "", nil)
}

// NewCollectionCommand creates an initialized command object
func NewCollectionCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewGetCommand(
		conf,
		ac,
		od.CollectionSchema,
		"collection",
		collectionsTransformation, "", nil)
}
