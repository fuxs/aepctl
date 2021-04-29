/*
Package cache consists of all caching relted functions and data structures.

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
package cache

import (
	"context"
	"time"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// SandboxCall encapsulates the sandbox list call
type SandboxCall struct {
	auth *api.AuthenticationConfig
}

// NewSandboxCall creates an initialized SandboxCall object
func NewSandboxCall(auth *api.AuthenticationConfig) APICall {
	return NewCachedAPICall(&SandboxCall{auth: auth})
}

// Call is the entry point
func (c *SandboxCall) Call() (interface{}, error) {
	return api.List(context.Background(), c.auth)
}

// NewSandboxCache creates an initilzed ListFileCache object for lists of sandboxes
func NewSandboxCache(auth *api.AuthenticationConfig, pp util.PathProvider) *ListFileCache {
	t := NewTransformList("sandboxes").V("name")
	return NewListFileCache(NewSandboxCall(auth), t, time.Hour*24, util.NewLazyPath(pp, "sandboxes.json"))
}
