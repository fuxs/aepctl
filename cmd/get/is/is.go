/*
Package is contains identity service related functions.

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
package is

import (
	"context"
	_ "embed"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

//go:embed trans/xid.yaml
var xidTransformation string

func NewXIDCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISParams{}
	cmd := &cobra.Command{
		Use:                   "xid (--namespace NAMESPACE_CODE|--ns-id NAMESPACE_ID) id",
		Short:                 "Display xid for id + namespace (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Aliases:               []string{"identity"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if pp.Namespace == "" && pp.NamespaceID == "" {
				pp.Namespace = "ECID"
			} else if pp.Namespace != "" && pp.NamespaceID != "" {
				helper.PrintError("Error: namespace code and namespace id specified. Please use either --namespace or --ns-id.", cmd)
			}
			helper.CheckErr(output.SetTransformationDesc(xidTransformation))
			pp.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.ISGetXID(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVarP(&pp.Namespace, "namespace", "n", "", "namespace code, ECID is default")
	flags.StringVar(&pp.NamespaceID, "ns-id", "", "namespace ID, e.g. 4 for ECID")
	flags.StringVar(&pp.Region, "region", "va7", "region for routing (va7 is default, nld2 is the alternative")
	return cmd
}

//go:embed trans/namespace.yaml
var nsTransformation string

func NewNamespaceCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	var imsOrg string
	cmd := &cobra.Command{
		Use:                   "namespace id",
		Short:                 "Display a namespace (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			helper.CheckErr(output.SetTransformationDesc(nsTransformation))
			if imsOrg == "" {
				helper.CheckErr(output.PrintResponse(api.ISGetNamespace(context.Background(), conf.Authentication, args[0])))
			} else {
				helper.CheckErr(output.PrintResponse(api.ISGetNamespaceIMSOrg(context.Background(), conf.Authentication, imsOrg, args[0])))
			}
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVar(&imsOrg, "ims-org", "", "IMS organization ID")
	return cmd
}

//go:embed trans/ids.yaml
var idsTransformation string

//go:embed trans/xids.yaml
var xidsTransformation string

func NewClusterCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISClusterParams{}
	cmd := &cobra.Command{
		Use:                   "id xid",
		Short:                 "Display a cluster of identities (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Aliases:               []string{"cluster"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if pp.Namespace != "" && pp.NamespaceID != "" {
				helper.PrintError("Error: namespace ids and codes are used together. Use either --namespace or --ns-id.", cmd)
			}
			if pp.Namespace == "" && pp.NamespaceID == "" {
				helper.CheckErr(output.SetTransformationDesc(xidsTransformation))
			} else {
				helper.CheckErr(output.SetTransformationDesc(idsTransformation))
			}
			pp.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.ISGetCluster(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVarP(&pp.Namespace, "namespace", "n", "", "namespace code, e.g. ECID")
	flags.StringVar(&pp.NamespaceID, "ns-id", "", "namespace ID, e.g. 4 for ECID")
	flags.StringVar(&pp.GraphType, "graph", "", "select identity graph. Private Graph is default, None for no graph")
	flags.StringVar(&pp.Region, "region", "va7", "region for routing (va7 is default, nld2 is the alternative")
	return cmd
}

func NewHistoryCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISClusterParams{}
	cmd := &cobra.Command{
		Use:                   "history xid",
		Short:                 "Display the history of a cluster (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if pp.Namespace != "" && pp.NamespaceID != "" {
				helper.PrintError("Error: namespace ids and codes are used together. Use either --namespace or --ns-id.", cmd)
			}
			/*if pp.Namespace == "" && pp.NamespaceID == "" {
				helper.CheckErr(output.SetTransformationDesc(xidsTransformation))
			} else {
				helper.CheckErr(output.SetTransformationDesc(idsTransformation))
			}*/
			pp.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.ISGetHistory(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVarP(&pp.Namespace, "namespace", "n", "", "namespace code, e.g. ECID")
	flags.StringVar(&pp.NamespaceID, "ns-id", "", "namespace ID, e.g. 4 for ECID")
	flags.StringVar(&pp.GraphType, "graph", "", "select identity graph. Private Graph is default, None for no graph")
	flags.StringVar(&pp.Region, "region", "va7", "region for routing (va7 is default, nld2 is the alternative")
	return cmd
}

func NewClustersCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISClustersParams{}
	cmd := &cobra.Command{
		Use:                   "ids xid+",
		Short:                 "Display multiple clusters of identities (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Aliases:               []string{"clusters"},
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			lns := len(pp.Namespaces)
			lid := len(pp.NamesapceIDs)
			largs := len(args)
			if lns > 0 && lid > 0 {
				helper.PrintError("Error: namespace ids and codes are used together. Use either --namespace or --ns-id.", cmd)
			}
			if lns+lid > largs {
				helper.PrintError("Error: more namespaces than ids", cmd)
			}
			if lns == 0 && lid == 0 {
				helper.CheckErr(output.SetTransformationDesc(xidsTransformation))
			} else {
				helper.CheckErr(output.SetTransformationDesc(idsTransformation))
			}
			pp.IDs = args
			helper.CheckErr(output.PrintResponse(api.ISGetClusters(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	//flags.StringArrayVarP(&pp.Namespaces, "namespace", "n", []string{}, "namespace code, e.g. ECID")
	flags.StringArrayVar(&pp.NamesapceIDs, "ns-id", []string{}, "namespace ID, e.g. 4 for ECID")
	flags.StringVar(&pp.GraphType, "graph", "", "select identity graph. Private Graph is default, None for no graph")
	flags.StringVar(&pp.Region, "region", "va7", "region for routing (va7 is default, nld2 is the alternative")
	return cmd
}

func NewHistoriesCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISClustersParams{}
	cmd := &cobra.Command{
		Use:                   "histories xid+",
		Short:                 "Display multiple histories of clusters (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			lns := len(pp.Namespaces)
			lid := len(pp.NamesapceIDs)
			largs := len(args)
			if lns > 0 && lid > 0 {
				helper.PrintError("Error: namespace ids and codes are used together. Use either --namespace or --ns-id.", cmd)
			}
			if lns+lid > largs {
				helper.PrintError("Error: more namespaces than ids", cmd)
			}
			/*if lns == 0 && lid == 0 {
				helper.CheckErr(output.SetTransformationDesc(xidsTransformation))
			} else {
				helper.CheckErr(output.SetTransformationDesc(idsTransformation))
			}*/
			pp.IDs = args
			helper.CheckErr(output.PrintResponse(api.ISGetHistories(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	//flags.StringArrayVarP(&pp.Namespaces, "namespace", "n", []string{}, "namespace code, e.g. ECID")
	flags.StringArrayVar(&pp.NamesapceIDs, "ns-id", []string{}, "namespace ID, e.g. 4 for ECID")
	flags.StringVar(&pp.GraphType, "graph", "", "select identity graph. Private Graph is default, None for no graph")
	flags.StringVar(&pp.Region, "region", "va7", "region for routing (va7 is default, nld2 is the alternative")
	return cmd
}

func NewMappingCommand(conf *helper.Configuration) *cobra.Command {
	output := &helper.OutputConf{}
	pp := &api.ISGetMappingParams{}
	cmd := &cobra.Command{
		Use:                   "mapping id",
		Short:                 "Display mapping for id (Identity Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd), output.ValidateFlags())
			if pp.Namespace != "" && pp.NamespaceID != "" {
				helper.PrintError("Error: namespace code and namespace id specified. Please use either --namespace or --ns-id.", cmd)
			}
			helper.CheckErr(output.SetTransformationDesc(nsTransformation))
			pp.ID = args[0]
			helper.CheckErr(output.PrintResponse(api.ISGetMapping(context.Background(), conf.Authentication, pp)))
		},
	}
	output.AddOutputFlags(cmd)
	flags := cmd.Flags()
	flags.StringVarP(&pp.Namespace, "namespace", "n", "", "namespace code, e.g. ECID")
	flags.StringVar(&pp.NamespaceID, "ns-id", "", "namespace ID, e.g. 4 for ECID")
	//flags.StringVar(&pp.TargetNS, "target", "", "target namespace ID, e.g. 4 for ECID")
	/*if err := cmd.MarkFlagRequired("target"); err != nil {
		helper.CheckErr(err)
		os.Exit(1)
	}*/
	return cmd
}
