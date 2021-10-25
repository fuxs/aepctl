/*
Package delete is the base for all delete commands.

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
package cancel

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewCancelQueryCommand creates an initialized command object
func NewCancelQueryCommand(conf *helper.Configuration) *cobra.Command {
	return NewCancelCommand(conf,
		api.QSCancelQuery,
		"query",
		"Cancel a query (Query Service)",
		"long",
		"example",
		"classes")
}
