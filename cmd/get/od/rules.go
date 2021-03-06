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

//go:embed trans/rules.yaml
var rulesTransformation string

// NewRulesCommand creates an initialized command object
func NewRulesCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewQueryCommand(
		conf,
		ac,
		od.RuleSchema,
		"rules",
		rulesTransformation, "", nil)
}

// NewRuleCommand creates an initialized command object
func NewRuleCommand(conf *helper.Configuration, ac *cache.AutoContainer) *cobra.Command {
	return NewGetCommand(
		conf,
		ac,
		od.RuleSchema,
		"rule",
		rulesTransformation, "", nil)
}
