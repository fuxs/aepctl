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
	"net/http"

	"github.com/fuxs/aepctl/util"
)

type SBGetSandboxParams string

func (p SBGetSandboxParams) Params() util.Params {
	return util.NewParams("name", string(p))
}

// SBGetSandbox returns the details of a sandbox with the given name
func SBGetSandbox(ctx context.Context, p *AuthenticationConfig, params SBGetSandboxParams) (*http.Response, error) {
	return SBGetSandboxP(ctx, p, params.Params())
}

// SBGetSandboxP returns the details of a sandbox with the given name
func SBGetSandboxP(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	name := params.GetForPath("name")
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/sandbox-management/sandboxes/%s", name)
}

// SBListAllSandboxes returns a list of all sandboxes
func SBListAllSandboxes(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return SBListAllSandboxesP(ctx, p, nil)
}

// SBListAllSandboxesP returns a list of all sandboxes. This variant implements the api.Func type.
func SBListAllSandboxesP(ctx context.Context, p *AuthenticationConfig, _ util.Params) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/sandbox-management/sandboxes")
}

// List returns a list of usable sandboxes
// TODO implement new generic API
func List(ctx context.Context, p *AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/foundation/sandbox-management/")
}

// SBListSandboxes returns a list of usable sandboxes
func SBListSandboxes(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return SBListSandboxesP(ctx, p, nil)
}

// SBListSandboxesP returns a list of usable sandboxes. This variant implements the api.Func type.
func SBListSandboxesP(ctx context.Context, p *AuthenticationConfig, _ util.Params) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/sandbox-management/")
}

// SBListSandboxTypes lists the available sandbox types
func SBListSandboxTypes(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return SBListSandboxTypesP(ctx, p, nil)
}

// SBListSandboxTypes lists the available sandbox types. This variant implements the api.Func type.
func SBListSandboxTypesP(ctx context.Context, p *AuthenticationConfig, _ util.Params) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/sandbox-management/sandboxTypes")
}
