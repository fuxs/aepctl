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

// NewDeleteActivitiesCommand creates an initialized command object
func NewDeleteClassCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.SRDeleteClass,
		"class",
		"Delete a class (Schema Registry)",
		"long",
		"example",
		"classes")
}

// NewDeleteDataTypeCommand creates an initialized command object
func NewDeleteDataTypeCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.SRDeleteDataType,
		"datatype",
		"Delete a data type (Schema Registry)",
		"long",
		"example",
		"datatypes",
		"data-type",
		"data-types")
}

// NewDeleteDataTypeCommand creates an initialized command object
func NewDeleteDescriptorCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.SRDeleteDescriptor,
		"descriptor",
		"Delete a descriptor (Schema Registry)",
		"long",
		"example",
		"descriptors")
}

// NewDeleteActivitiesCommand creates an initialized command object
func NewDeleteFieldGroupCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.SRDeleteFieldGroup,
		"fieldgroup",
		"Delete a field group (Schema Registry)",
		"long",
		"example",
		"field-group",
		"fieldgroups",
		"field-groups")
}

// NewDeleteActivitiesCommand creates an initialized command object
func NewDeleteSchemaCommand(conf *helper.Configuration) *cobra.Command {
	return NewDeleteCommand(conf,
		api.SRDeleteSchema,
		"schema",
		"Delete a schema (Schema Registry)",
		"long",
		"example",
		"schemas")
}
