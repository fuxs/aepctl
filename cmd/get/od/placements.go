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
	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type placementTransformer struct{}

func (*placementTransformer) Header(wide bool) []string {
	if wide {
		return []string{"ID", "NAME", "CHANNEL TYPE", "CONTENT TYPE", "LAST MODIFIED", "DESCRIPTION"}
	}
	return []string{"NAME", "CHANNEL TYPE", "CONTENT TYPE", "LAST MODIFIED", "DESCRIPTION"}
}

func (*placementTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*placementTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	if wide {
		return w.Write(
			s.Str("@id"),
			s.Str("xdm:name"),
			helper.ChannelLToS.Lookup(s.Str("xdm:channel")),
			helper.ContentLToS.Lookup(s.Str("xdm:componentType")),
			util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
			s.Str("xdm:description"))
	}
	return w.Write(
		s.Str("xdm:name"),
		helper.ChannelLToS.Lookup(s.Str("xdm:channel")),
		helper.ContentLToS.Lookup(s.Str("xdm:componentType")),
		util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
		s.Str("xdm:description"))
}

func (*placementTransformer) Iterator(*util.JSONCursor) (util.JSONResponse, error) {
	return nil, nil
}

// NewPlacementsCommand creates an initialized command object
func NewPlacementsCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	pt := &placementTransformer{}
	return NewQueryCommand(
		conf,
		ac,
		od.PlacementSchema,
		"placements",
		pt)
}

// NewPlacementCommand creates an initialized command object
func NewPlacementCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	pt := &placementTransformer{}
	return NewGetCommand(
		conf,
		ac,
		od.PlacementSchema,
		"placement",
		pt)
}
