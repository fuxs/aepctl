/*
Package create is the base for all create commands.

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
package create

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

// NewClassCommand creates an initialized command object
func NewClassCommand(conf *helper.Configuration) *cobra.Command {
	return NewCreateCommand(conf,
		api.SRCreateClass,
		"class",
		"Create a class (Schema Registry)",
		"long",
		"example",
		"classes",
	)
}

// NewDataTypeCommand creates an initialized command object
func NewDataTypeCommand(conf *helper.Configuration) *cobra.Command {
	return NewCreateCommand(conf,
		api.SRCreateDataType,
		"datatype",
		"Create a data type (Schema Registry)",
		"long",
		"example",
		"datatypes",
		"data-type",
		"data-types",
	)
}

// NewDescriptorCommand creates an initialized command object
func NewDescriptorCommand(conf *helper.Configuration) *cobra.Command {
	return NewCreateCommand(conf,
		api.SRCreateDescriptor,
		"descriptor",
		"Create a descriptor (Schema Registry)",
		"long",
		"example",
		"descriptors",
	)
}

// NewFieldGroupCommand creates an initialized command object
func NewFieldGroupCommand(conf *helper.Configuration) *cobra.Command {
	return NewCreateCommand(conf,
		api.SRCreateFieldGroup,
		"fieldgroup",
		"Create a fieldgroup (Schema Registry)",
		"long",
		"example",
		"fieldgroups",
	)
}

// NewSchemaCommand creates an initialized command object
func NewSchemaCommand(conf *helper.Configuration) *cobra.Command {
	return NewCreateCommand(conf,
		api.SRCreateSchema,
		"schema",
		"Create a schema (Schema Registry)",
		"long",
		"example",
		"schemas",
	)
}
