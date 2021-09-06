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

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/api/od"
)

// APICall is a simple interface for all simplified REST API calls
type APICall interface {
	Call() (interface{}, error)
}

// CachedAPICall just stores a copy of the response if call was successful.
// Implements APICall interface.
type CachedAPICall struct {
	apiCall APICall
	obj     interface{}
}

// NewCachedAPICall creates an initialized CachedAPICall object
func NewCachedAPICall(apiCall APICall) APICall {
	return &CachedAPICall{apiCall: apiCall}
}

// Call is the entry point
func (c *CachedAPICall) Call() (interface{}, error) {
	if c.obj != nil {
		return c.obj, nil
	}
	obj, err := c.apiCall.Call()
	if err != nil {
		return nil, err
	}
	c.obj = obj
	return obj, nil

}

// ContainerCall encapsulates the call of the container service
type ContainerCall struct {
	auth *api.AuthenticationConfig
}

// NewContainerCall creates an initialized ContainerCall object
func NewContainerCall(auth *api.AuthenticationConfig) APICall {
	return NewCachedAPICall(&ContainerCall{auth: auth})
}

// Call is the entry point
func (c *ContainerCall) Call() (interface{}, error) {
	return od.ListContainer(context.Background(), c.auth)
}

// ODCall is a generic encapsulation for all offer decisioning calls
type ODCall struct {
	ac     *AutoContainer
	schema string
}

//NewODCall creates an initialized ODCall object
func NewODCall(ac *AutoContainer, schema string) APICall {
	return NewCachedAPICall(&ODCall{ac: ac, schema: schema})
}

// Call is the entry point
func (c *ODCall) Call() (interface{}, error) {
	cid, err := c.ac.Get()
	if err != nil {
		return nil, err
	}
	param := &od.ListParam{
		ContainerID: cid,
		Schema:      c.schema,
	}
	return od.List(context.Background(), c.ac.Auth, param)
}
