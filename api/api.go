/*
Package api is the base for all aep rest functions.

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
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fuxs/aepctl/util"
)

type FuncID func(context.Context, *AuthenticationConfig, string) (*http.Response, error)
type Func func(context.Context, *AuthenticationConfig, util.Params) (*http.Response, error)
type FuncTable map[string]Func
type Parameters interface {
	Params() util.Params
}

type ParametersE interface {
	Params() (util.Params, error)
}

var All = map[string]Func{
	"ACGetPermissionsAndResources": ACGetPermissionsAndResourcesP,
	"ACGetEffecticeACLPolicies":    ACGetEffecticeACLPoliciesP,
	"CatalogGetBatches":            CatalogGetBatchesP,
	"CatalogGetDatasets":           CatalogGetDatasetsP,
	"DADownload":                   DADownloadP,
	"DAGetFiles":                   DAGetFilesP,
	"DAGetFils":                    DAGetFileP,
	"ODGet":                        ODGetP,
	"ODQuery":                      ODQueryP,
	"FlowGetConnections":           FlowGetConnectionsP,
	"SBGetSandbox":                 SBGetSandboxP,
	"SBListAllSandboxes":           SBListAllSandboxesP,
	"SBListSandboxes":              SBListSandboxesP,
	"SBListSandboxTypes":           SBListSandboxTypesP,
	"SRGetBehaviorsP":              SRListBehaviorsP,
	"SRGetBehaviorP":               SRGetBehaviorP,
	"SRGetStats":                   SRGetStatsP,
	"SRGetSchemas":                 SRListSchemasP,
	"SRGetSchema":                  SRGetSchemaP,
	"UPSGetEntities":               UPSGetEntitiesP,
}

func (ft FuncTable) Call(name string, ctx context.Context, auth *AuthenticationConfig, params util.Params) (*http.Response, error) {
	f := ft[name]
	if f == nil {
		return nil, fmt.Errorf("function %s not found", name)
	}
	return f(ctx, auth, params)
}
