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
package ups

import (
	"context"
	_ "embed"
	"time"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/profile.yaml
var profileTransformation string

// NewEntitiesCommand creates an initialized command object
func NewEntitiesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.UPSEntitiesParams{}
	cmd := &cobra.Command{
		Use:                   "entities",
		Short:                 "Display all entities",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationDesc(profileTransformation))
			output.StreamResultRaw(api.UPSGetEntities(context.Background(), conf.Authentication, p))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVar(&p.TimeFormat, "time-format", time.RFC3339, "format for date parsing, default is '2006-01-02T15:04:05Z07:00' (RFC3339)")
	flags.StringVar(&p.Schema, "schema", "_xdm.context.profile", "XED schema class name, default is _xdm.context.profile")
	flags.StringVar(&p.RelatedSchema, "related-schema", "", "Must be set if schema is _xdm.context.experienceevent")
	flags.StringVarP(&p.ID, "id", "i", "", "ID of the entity. For Native XID lookup, use <XID> and leave ns empty")
	flags.StringVarP(&p.NS, "ns", "n", "", "Identity namespace code")
	flags.StringVar(&p.RelatedID, "related-id", "", "ID of the entity that the ExperienceEvents are associated with.")
	flags.StringVar(&p.RelatedNS, "related-ns", "", "Identity namespace code")
	flags.StringVar(&p.Fields, "fields", "", "Fields for the model object. By default, all fields will be fetched")
	flags.StringVar(&p.MP, "mp", "", "ID of the merge policy")
	flags.StringVar(&p.Start, "start", "", "Start time of Time range filter for ExperienceEvents")
	flags.StringVar(&p.End, "end", "", "End time of Time range filter for ExperienceEvents")
	flags.IntVar(&p.Limit, "limit", 0, "Number of records to return from the result")
	flags.StringVar(&p.Order, "order", "", "The sort order of retrieved ExperienceEvents by timestamp")
	flags.StringVar(&p.Property, "property", "", "End time of Time range filter for ExperienceEvents")
	flags.BoolVar(&p.CA, "ca", false, "Feature flag for enabling computed attributes for lookup")
	return cmd
}

func NewProfileCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.UPSEntitiesParams{Schema: "_xdm.context.profile"}
	cmd := &cobra.Command{
		Use:                   "profile",
		Short:                 "Display a profile",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationDesc(profileTransformation))
			p.ID = args[0]
			output.StreamResultRaw(api.UPSGetEntities(context.Background(), conf.Authentication, p))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	// flags.StringVar(&p.Schema, "schema", "_xdm.context.profile", "XED schema class name, default is _xdm.context.profile")
	// flags.StringVar(&p.RelatedSchema, "related-schema", "", "Must be set if schema is _xdm.context.experienceevent")
	// flags.StringVarP(&p.ID, "id", "i", "", "ID of the entity. For Native XID lookup, use <XID> and leave ns empty")
	flags.StringVarP(&p.NS, "ns", "n", "", "Identity namespace code")
	// flags.StringVar(&p.RelatedID, "related-id", "", "ID of the entity that the ExperienceEvents are associated with.")
	// flags.StringVar(&p.RelatedNS, "related-ns", "", "Identity namespace code")
	flags.StringVar(&p.Fields, "fields", "", "Fields for the model object. By default, all fields will be fetched")
	flags.StringVar(&p.MP, "mp", "", "ID of the merge policy")
	// flags.StringVar(&p.Start, "start", "", "Start time of Time range filter for ExperienceEvents")
	// flags.StringVar(&p.End, "end", "", "End time of Time range filter for ExperienceEvents")
	// flags.IntVar(&p.Limit, "limit", 0, "Number of records to return from the result")
	// flags.StringVar(&p.Order, "order", "", "The sort order of retrieved ExperienceEvents by timestamp")
	// flags.StringVar(&p.Property, "property", "", "End time of Time range filter for ExperienceEvents")
	flags.BoolVar(&p.CA, "ca", false, "Feature flag for enabling computed attributes for lookup")
	return cmd
}

func NewEventsCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	p := &api.UPSEntitiesParams{Schema: "_xdm.context.experienceevent", RelatedSchema: "_xdm.context.profile"}
	cmd := &cobra.Command{
		Use:                   "events",
		Short:                 "Display events for a profile",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErrs(output.SetTransformationDesc(profileTransformation))
			p.RelatedID = args[0]
			output.StreamResultRaw(api.UPSGetEntities(context.Background(), conf.Authentication, p))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	// flags.StringVar(&p.Schema, "schema", "_xdm.context.profile", "XED schema class name, default is _xdm.context.profile")
	// flags.StringVar(&p.RelatedSchema, "related-schema", "", "Must be set if schema is _xdm.context.experienceevent")
	// flags.StringVarP(&p.ID, "id", "i", "", "ID of the entity. For Native XID lookup, use <XID> and leave ns empty")
	//flags.StringVarP(&p.NS, "ns", "n", "", "Identity namespace code")
	// flags.StringVar(&p.RelatedID, "related-id", "", "ID of the entity that the ExperienceEvents are associated with.")
	flags.StringVar(&p.RelatedNS, "related-ns", "", "Identity namespace code")
	flags.StringVar(&p.Fields, "fields", "", "Fields for the model object. By default, all fields will be fetched")
	flags.StringVar(&p.MP, "mp", "", "ID of the merge policy")
	flags.StringVar(&p.Start, "start", "", "Start time of Time range filter for ExperienceEvents")
	flags.StringVar(&p.End, "end", "", "End time of Time range filter for ExperienceEvents")
	flags.IntVar(&p.Limit, "limit", 0, "Number of records to return from the result")
	flags.StringVar(&p.Order, "order", "", "The sort order of retrieved ExperienceEvents by timestamp")
	flags.StringVar(&p.Property, "property", "", "End time of Time range filter for ExperienceEvents")
	flags.BoolVar(&p.CA, "ca", false, "Feature flag for enabling computed attributes for lookup")
	return cmd
}
