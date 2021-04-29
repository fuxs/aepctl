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
	"strconv"
	"strings"

	"github.com/fuxs/aepctl/util"
)

// GetStats returns schema registry informations
func SRGetStats(ctx context.Context, p *AuthenticationConfig) (*http.Response, error) {
	return SRGetStatsP(ctx, p, nil)
}

// GetStatsP returns schema registry informations
func SRGetStatsP(ctx context.Context, p *AuthenticationConfig, _ util.Params) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats")
}

type SRGetSchemasParams struct {
	Properties string
	OrderBy    string
	Start      string
	Limit      uint
	Global     bool
	Full       bool
}

func (p *SRGetSchemasParams) Params() util.Params {
	var (
		limit  string
		cid    string
		accept string
	)
	if p.Limit > 0 {
		limit = strconv.FormatUint(uint64(p.Limit), 10)
	}
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	if p.Full {
		accept = "application/vnd.adobe.xed+json"
	} else {
		accept = "application/vnd.adobe.xed-id+json"
	}
	return util.NewParams(
		"properties", p.Properties,
		"orderby", p.OrderBy,
		"start", p.Start,
		"limit", limit,
		"-cid", cid,
		"-accept", accept,
	)
}

func SRGetSchemas(ctx context.Context, a *AuthenticationConfig, p *SRGetSchemasParams) (*http.Response, error) {
	return SRGetSchemasP(ctx, a, p.Params())
}

func SRGetSchemasP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	header := map[string]string{"Accept": p.Get("-accept")}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/schemas%s", p.Get("-cid"), p.EncodeWithout("-cid", "-accept"))
}

type SRGetSchemaParams struct {
	ID          string
	Version     string
	Global      bool
	Full        bool
	NoText      bool
	Descriptors bool
}

func (p *SRGetSchemaParams) Params() util.Params {
	var (
		cid string
		sb  strings.Builder
	)
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	sb.WriteString("application/vnd.adobe.xed")
	if p.Full {
		sb.WriteString("-full")
	}
	if p.NoText {
		sb.WriteString("-notext")
	} else if p.Full && p.Descriptors {
		sb.WriteString("-desc")
	}
	sb.WriteString("+json; version=")
	if p.Version != "" {
		sb.WriteString(p.Version)
	} else {
		sb.WriteString("1")
	}

	return util.NewParams(
		"id", p.ID,
		"cid", cid,
		"accept", sb.String(),
	)
}

func SRGetSchema(ctx context.Context, a *AuthenticationConfig, p *SRGetSchemaParams) (*http.Response, error) {
	return SRGetSchemaP(ctx, a, p.Params())
}

func SRGetSchemaP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	header := map[string]string{"Accept": p.Get("accept")}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/schemas/%s", p.GetForPath("cid"), p.GetForPath("id"))
}
