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
package delete

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewDeleteQueryCommand creates an initialized command object
func NewDeleteQueryCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.QSDeleteQuery,
		"query",
		"Delete a query (Query Service)",
		"long",
		"example",
		"classes")
}

// NewDeleteScheduledQueryCommand creates an initialized command object
func NewDeleteScheduledQueryCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.QSDeleteSchedule,
		"schedule",
		"Delete a scheduled query (Query Service)",
		"long",
		"example",
		"schedules")
}

// NewDeleteQueryTemplateCommand creates an initialized command object
func NewDeleteQueryTemplateCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.QSDeleteQueryTemplate,
		"template",
		"Delete a query template (Query Service)",
		"long",
		"example",
		"templates")
}
