/*
Package acl consists of access control list functions.

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
package acl

import (
	"context"

	"github.com/fuxs/aepctl/api"
)

// GetPermissionsAndResources returns the access control policies.
func GetPermissionsAndResources(ctx context.Context, p *api.AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx, "https://platform.adobe.io/data/foundation/access-control/acl/reference")
}

// GetEffecticeACLPolicies returns the effective acl policies
func GetEffecticeACLPolicies(ctx context.Context, p *api.AuthenticationConfig, urls []string) (interface{}, error) {
	return p.PostJSONRequest(ctx, urls, "https://platform.adobe.io/data/foundation/access-control/acl/effective-policies")
}
