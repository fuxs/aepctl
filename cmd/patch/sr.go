/*
Package patch contains patch command related functions.

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
package patch

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewClassCommand creates an initialized command object
func NewClassCommand(conf *helper.Configuration) *cobra.Command {
	return NewPatchCommand(conf,
		api.SRPatchClass,
		"class",
		"Update a class (Schema Registry)",
		"long",
		"example",
		"classes",
	)
}

// NewDataTypeCommand creates an initialized command object
func NewDataTypeCommand(conf *helper.Configuration) *cobra.Command {
	return NewPatchCommand(conf,
		api.SRPatchDataType,
		"datatype",
		"Update a data type (Schema Registry)",
		"long",
		"example",
		"datatypes",
		"data-type",
		"data-types",
	)
}

// NewFieldGroupCommand creates an initialized command object
func NewFieldGroupCommand(conf *helper.Configuration) *cobra.Command {
	return NewPatchCommand(conf,
		api.SRPatchFieldGroup,
		"fieldgroup",
		"Update a fieldgroup (Schema Registry)",
		"long",
		"example",
		"fieldgroups",
	)
}

// NewSchemaCommand creates an initialized command object
func NewSchemaCommand(conf *helper.Configuration) *cobra.Command {
	return NewPatchCommand(conf,
		api.SRPatchSchema,
		"schema",
		"Update a schema (Schema Registry)",
		"long",
		"example",
		"schemas",
	)
}
