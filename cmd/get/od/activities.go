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

type activityTransformer struct {
	idc *cache.MapMemCache
}

func newActivityTransformer(ac *cache.AutoContainer) *activityTransformer {
	// get list of placements and store map[@id]channel
	t := cache.NewODTrans().K("_instance", "@id").V("_instance", "xdm:channel")
	c := cache.NewODCall(ac, od.PlacementSchema)
	idc := cache.NewMapMemCache(c, t)
	return &activityTransformer{
		idc: idc,
	}
}

func (*activityTransformer) Header(wide bool) []string {
	return []string{"NAME", "STATUS", "START DATE", "END DATE", "CHANNEL TYPE", "LAST MODIFIED"}
}

func (*activityTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (t *activityTransformer) WriteRow(q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	return w.Write(
		s.Str("xdm:name"),
		StatusMapper.Lookup(s.Str("xdm:status")),
		util.LocalTimeStrCustom(s.Str("xdm:startDate"), shortDate),
		util.LocalTimeStrCustom(s.Str("xdm:endDate"), shortDate),
		s.Path("xdm:criteria").Concat(",", func(q *util.Query) string {
			id := q.Path("xdm:placements").Get(0).String()
			return helper.ChannelLToS.Lookup(t.idc.Lookup(id))
		}),
		util.LocalTimeStrCustom(q.Str("repo:lastModifiedDate"), longDate),
	)
}

func (*activityTransformer) Iterator(*util.JSONCursor) (util.JSONResponse, error) {
	return nil, nil
}

// NewActivitiesCommand creates an initialized command object
func NewActivitiesCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	at := newActivityTransformer(ac)
	return NewQueryCommand(
		conf,
		ac,
		od.ActivitySchema,
		"activities",
		at)
}

// NewActivityCommand creates an initialized command object
func NewActivityCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	at := newActivityTransformer(ac)
	return NewGetCommand(
		conf,
		ac,
		od.ActivitySchema,
		"activity",
		at)
}
