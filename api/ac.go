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

// ACGetPermissionsAndResources returns the access control policies.
func ACGetPermissionsAndResources(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return ACGetPermissionsAndResourcesP(ctx, p, nil)
}

// ACGetPermissionsAndResources returns the access control policies.
func ACGetPermissionsAndResourcesP(ctx context.Context, p *AuthenticationConfig, _ util.Params) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/access-control/acl/reference")
}

type ACGetEffecticeACLPoliciesParams []string

func (p ACGetEffecticeACLPoliciesParams) Params() util.Params {
	return util.Params{"urls": []string(p)}
}

// ACGetEffecticeACLPolicies returns the effective acl policies
func ACGetEffecticeACLPolicies(ctx context.Context, p *AuthenticationConfig, params ACGetEffecticeACLPoliciesParams) (*http.Response, error) {
	return ACGetEffecticeACLPoliciesP(ctx, p, params.Params())
}

// ACGetEffecticeACLPolicies returns the effective acl policies
func ACGetEffecticeACLPoliciesP(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	urls := params["urls"]
	return p.PostJSONRequestRaw(ctx, urls, "https://platform.adobe.io/data/foundation/access-control/acl/effective-policies")
}
