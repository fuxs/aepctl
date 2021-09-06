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
	"errors"
	"net/http"
	"strconv"
)

// ODQueryParames contains all available parameters for offer decisioning
// queries.
type ODQueryParames struct {
	ContainerID string
	Schema      string
	Query       string
	QOP         string
	Field       string
	OrderBy     string
	Limit       int
}

// Params returns the parameters in generic Request format
func (p *ODQueryParames) Request() *Request {
	var limit string
	if p.Limit > 0 {
		limit = strconv.FormatInt(int64(p.Limit), 10)
	}
	req := NewRequest(
		"schema", p.Schema,
		"q", p.Query,
		"qop", p.QOP,
		"field", p.Field,
		"oderBy", p.OrderBy,
		"limit", limit,
	)
	req.SetValue("containerID", p.ContainerID)
	return req
}

// ODQueryP sends a query to the offer decisioning API
func ODQueryP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	containerID := params.GetValuePath("containerID")
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	return p.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/queries/core/search%s",
		containerID,
		params.EncodedQuery(),
	)
}

type ODGetParams struct {
	ID          string
	ContainerID string
	Schema      string
}

// Params returns the parameters in generic Request format
func (p *ODGetParams) Request() *Request {
	req := NewRequest("schema", p.Schema, "id", p.ID)
	req.SetValue("containerID", p.ContainerID)
	return req
}

// ODGet returns a collection by name (wild cards are supported) or id (exact match)
func ODGet(ctx context.Context, p *AuthenticationConfig, params *ODGetParams) (*http.Response, error) {
	return ODGetP(ctx, p, params.Request())
}

// ODGetP returns a collection by name (wild cards are supported) or id (exact match)
func ODGetP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	containerID := params.GetValuePath("containerID")
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	return p.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/instances%s",
		containerID,
		params.EncodedQuery(),
	)
}
