/*
Package get contains get command related functions.

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
package get

import (
	"github.com/fuxs/aepctl/cmd/get/sr"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewSRCommand creates an initialized command object
func NewSRCommand(conf *helper.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sr",
		Short: "Schema Registry",
	}
	cmd.AddCommand(sr.NewStatsCommand(conf))
	cmd.AddCommand(sr.NewSchemasCommand(conf))
	cmd.AddCommand(sr.NewSchemaCommand(conf))
	return cmd
}
