/*
Package token contains all token related functions.

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
	"strconv"

	"github.com/fuxs/aepctl/util"
)

// GetStats returns schema registry informations
func SRGetStatsRaw(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats")
}

// SRGetStats returns schema registry informations
func SRGetStats(ctx context.Context, p *AuthenticationConfig) (util.JSONResponse, error) {
	return NewJSONIterator(p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats"))
}

type SRGetSchemasParam struct {
	Properties string
	OrderBy    string
	Start      string
	Limit      uint
	Global     bool
	Full       bool
}

func (p *SRGetSchemasParam) toParams() string {
	var limit string
	if p.Limit > 0 {
		limit = strconv.FormatUint(uint64(p.Limit), 10)
	}
	return util.Par("properties", p.Properties, "orderby", p.OrderBy, "start", p.Start, "limit", limit)
}

func SRGetSchemasRaw(ctx context.Context, a *AuthenticationConfig, p *SRGetSchemasParam) (*http.Response, error) {
	var cid string
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	accept := "application/vnd.adobe.xed-id+json"
	if p.Full {
		accept = "application/vnd.adobe.xed+json"
	}
	header := map[string]string{"Accept": accept}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/schemas%s", cid, p.toParams())
}

func SRGetSchemaRaw(ctx context.Context, a *AuthenticationConfig, p *SRGetSchemasParam, id string) (*http.Response, error) {
	var cid string
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	accept := "application/vnd.adobe.xed-id+json"
	if p.Full {
		accept = "application/vnd.adobe.xed+json"
	}
	header := map[string]string{"Accept": accept}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/schemas/%s", cid, id)
}
