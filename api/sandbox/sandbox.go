/*
Package sandbox contains all sandbox related functions.

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
package sandbox

import (
	"context"

	"github.com/fuxs/aepctl/api"
)

// Get returns the details of a sandbox with the given name
func Get(ctx context.Context, p *api.AuthenticationConfig, name string) (interface{}, error) {
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/foundation/sandbox-management/sandboxes/%s", name)
}

// ListAll returns a list of all sandboxes
func ListAll(ctx context.Context, p *api.AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/foundation/sandbox-management/sandboxes")
}

// List returns a list of usable sandboxes
func List(ctx context.Context, p *api.AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/foundation/sandbox-management/")
}

// ListTypes lists the available sandbox types
func ListTypes(ctx context.Context, p *api.AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/foundation/sandbox-management/sandboxTypes")
}
