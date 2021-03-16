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
	"io"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/cache"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

type ruleTransformer struct{}

func (*ruleTransformer) Header(wide bool) []string {
	return []string{"NAME", "LAST MODIFIED", "DESCRIPTION"}
}

func (*ruleTransformer) Preprocess(i util.JSONResponse) error {
	if err := i.Path("_embedded", "results"); err != nil {
		return err
	}
	return i.EnterArray()
}

func (*ruleTransformer) WriteRow(name string, q *util.Query, w *util.RowWriter, wide bool) error {
	s := q.Path("_instance")
	return w.Write(
		s.Str("xdm:name"),
		util.LocalTimeStr(q.Str("repo:lastModifiedDate")),
		s.Str("xdm:description"),
	)
}

func (*ruleTransformer) Iterator(io.ReadCloser) (util.JSONResponse, error) {
	return nil, nil
}

// NewRulesCommand creates an initialized command object
func NewRulesCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	rt := &ruleTransformer{}
	return NewQueryCommand(
		conf,
		ac,
		od.RuleSchema,
		"rules",
		rt)
}

// NewRuleCommand creates an initialized command object
func NewRuleCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	rt := &ruleTransformer{}
	return NewGetCommand(
		conf,
		ac,
		od.RuleSchema,
		"rule",
		rt)
}
