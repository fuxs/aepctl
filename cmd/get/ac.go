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
	"context"

	"github.com/fuxs/aepctl/api/acl"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

var (
	// TODO long description
	aclLong = util.LongDesc(`
	Display all or effictive permissions.
	
	`)

	// TODO provide more examples
	aclExample = util.Example(`
	# Get all permissions
	aepctl get acl

	# Get specific resource types and permissions
	aepctl get acl /resource-types/activations /permissions/activate-destinations
	`)
)

type acTransformer struct{}

func (t *acTransformer) ToTable(i interface{}) (*util.Table, error) {
	table := util.NewTable([]string{"OPERATION", "OBJECT", "R", "W", "D"}, 64)
	util.NewQuery(i).Path("permissions").RangeAttributes(func(s string, q *util.Query) {
		operationName := s
		q.RangeAttributes(func(object string, q *util.Query) {
			permissions := q.Strings()
			read := util.ContainsS("read", permissions)
			write := util.ContainsS("write", permissions)
			delete := util.ContainsS("delete", permissions)
			table.Append(map[string]interface{}{
				"OPERATION": operationName,
				"OBJECT":    object,
				"R":         read,
				"W":         write,
				"D":         delete,
			})
			operationName = ""
		})
	})
	return table, nil
}

func (t *acTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	table := util.NewTable([]string{"OPERATION", "OBJECT", "READ", "WRITE", "DELETE"}, 64)
	util.NewQuery(i).Path("permissions").RangeAttributes(func(s string, q *util.Query) {
		operationName := s
		q.RangeAttributes(func(object string, q *util.Query) {
			permissions := q.Strings()
			read := util.ContainsS("read", permissions)
			write := util.ContainsS("write", permissions)
			delete := util.ContainsS("delete", permissions)
			table.Append(map[string]interface{}{
				"OPERATION": operationName,
				"OBJECT":    object,
				"READ":      read,
				"WRITE":     write,
				"DELETE":    delete,
			})
			operationName = ""
		})
	})
	return table, nil
}

type effectiveTransformer struct{}

func (t *effectiveTransformer) ToTable(i interface{}) (*util.Table, error) {
	table := util.NewTable([]string{"OBJECT", "VALUES"}, 64)
	util.NewQuery(i).Path("policies").RangeAttributes(func(s string, q *util.Query) {
		table.Append(map[string]interface{}{
			"OBJECT": s,
			"VALUES": q.Concat(",", func(q *util.Query) string { return q.String() }),
		})
	})
	return table, nil
}

func (t *effectiveTransformer) ToWideTable(i interface{}) (*util.Table, error) {
	return t.ToTable(i)
}

var validArgs = []string{
	"/permissions/activate-destinations",
	"/permissions/evaluate-segments",
	"/permissions/execute-decisioning-activities",
	"/permissions/export-audience-for-segment",
	"/permissions/manage-datasets",
	"/permissions/manage-decisioning-activities",
	"/permissions/manage-decisioning-options",
	"/permissions/manage-destinations",
	"/permissions/manage-dsw",
	"/permissions/manage-dule-labels",
	"/permissions/manage-dule-policies",
	"/permissions/manage-identity-namespaces",
	"/permissions/manage-privacy-workflows",
	"/permissions/manage-profile-configs",
	"/permissions/manage-profiles",
	"/permissions/manage-queries",
	"/permissions/manage-schemas",
	"/permissions/manage-segments",
	"/permissions/manage-sources",
	"/permissions/reset-sandboxes",
	"/permissions/view-datasets",
	"/permissions/view-destinations",
	"/permissions/view-dule-labels",
	"/permissions/view-dule-policies",
	"/permissions/view-identity-namespaces",
	"/permissions/view-monitoring-dashboard",
	"/permissions/view-privacy-workflows",
	"/permissions/view-profile-configs",
	"/permissions/view-profiles",
	"/permissions/view-sandboxes",
	"/permissions/view-schemas",
	"/permissions/view-segments",
	"/permissions/view-sources",
	"/resource-types/activation-associations",
	"/resource-types/activations",
	"/resource-types/activities",
	"/resource-types/analytics-source",
	"/resource-types/audience-manager-source",
	"/resource-types/bizible-source",
	"/resource-types/connection",
	"/resource-types/customer-attributes-source",
	"/resource-types/data-science-workspace",
	"/resource-types/dataset-preview",
	"/resource-types/datasets",
	"/resource-types/dule-label",
	"/resource-types/dule-policy",
	"/resource-types/enterprise-source",
	"/resource-types/identity-descriptor",
	"/resource-types/identity-namespaces",
	"/resource-types/launch-source",
	"/resource-types/marketing-action",
	"/resource-types/marketo-source",
	"/resource-types/monitoring",
	"/resource-types/offers",
	"/resource-types/placements",
	"/resource-types/privacy-consent",
	"/resource-types/privacy-content-delivery",
	"/resource-types/privacy-job",
	"/resource-types/profile-configs",
	"/resource-types/profile-datasets",
	"/resource-types/profiles",
	"/resource-types/query",
	"/resource-types/relationship-descriptor",
	"/resource-types/sandboxes",
	"/resource-types/schemas",
	"/resource-types/segment-jobs",
	"/resource-types/segments",
	"/resource-types/streaming-source",
}

// NewACCommand creates an initialized command object
func NewACCommand(auth *helper.Authentication) *cobra.Command {
	output := helper.NewOutputConf(&acTransformer{})
	cmd := &cobra.Command{
		Use:                   "ac [(RESOURCE | PERMISSION)*]",
		Short:                 "Display all or effictive permissions",
		Long:                  aclLong,
		Example:               aclExample,
		DisableFlagsInUseLine: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return util.Difference(validArgs, args), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErr(output.ValidateFlags())
			ctx := context.Background()
			if len(args) == 0 {
				output.PrintResult(acl.GetPermissionsAndResources(ctx, auth.Config))
			} else {
				output.SetTransformation(&effectiveTransformer{})
				output.PrintResult(acl.GetEffecticeACLPolicies(ctx, auth.Config, args))
			}
		},
	}
	output.AddOutputFlags(cmd)
	return cmd
}
