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
	"io"
	"net/http"
	"net/url"
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

type SRFormat struct {
	Short       bool
	Full        bool
	NoText      bool
	Descriptors bool
	Version     string
}

// accept builds the accept string
func (s *SRFormat) Accept() string {
	if s.Short {
		return "application/vnd.adobe.xed-id+json"
	}
	var sb strings.Builder
	sb.WriteString("application/vnd.adobe.xed")
	if s.Full {
		sb.WriteString("-full")
	}
	if s.Descriptors {
		sb.WriteString("-desc")
	}
	if s.NoText {
		sb.WriteString("-notext")
	}

	sb.WriteString("+json")
	if s.Version != "" {
		sb.WriteString("; version=")
		sb.WriteString(s.Version)
	}

	return sb.String()
}

func (s *SRFormat) Params() util.Params {
	return util.NewParams(
		"accept", s.Accept(),
	)
}

const (
	// JSONOut is used for JSON
	AcceptIDs = iota
	AcceptLinks
	AcceptObjects
)

type SRDescriptorFormat int

func (f SRDescriptorFormat) Accept() string {
	switch f {
	case AcceptIDs:
		return "application/vnd.adobe.xdm-v2-id+json"
	case AcceptLinks:
		return "application/vnd.adobe.xdm-v2-link+json"
	default:
		return "application/vnd.adobe.xdm-v2+json"
	}
}

type SRBaseParams struct {
	Property string
	OrderBy  string
	Start    string
	Limit    uint
	Global   bool
}

func (p *SRBaseParams) Params() util.Params {
	var (
		limit string
		cid   string
	)
	if p.Limit > 0 {
		limit = strconv.FormatUint(uint64(p.Limit), 10)
	}
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	return util.NewParams(
		"property", p.Property,
		"orderby", p.OrderBy,
		"start", p.Start,
		"limit", limit,
		"-cid", cid,
	)
}

type SRListParams struct {
	SRBaseParams
	SRFormat
}

func (p *SRListParams) Params() util.Params {
	result := p.SRBaseParams.Params()
	result.Add("-accept", p.Accept())
	return result
}

type SRListDescriptorsParams struct {
	SRBaseParams
	SRDescriptorFormat
}

func (p *SRListDescriptorsParams) Params() util.Params {
	result := p.SRBaseParams.Params()
	result.Add("-accept", p.Accept())
	return result
}

func srList(ctx context.Context, a *AuthenticationConfig, p util.Params, res string) (*http.Response, error) {
	header := map[string]string{"Accept": p.Get("-accept")}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/%s%s", p.Get("-cid"), res, p.EncodeWithout("-cid", "-accept"))
}

func SRListBehaviors(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "behaviors")
}

func SRListBehaviorsP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "behaviors")
}

func SRListClasses(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "classes")
}

func SRListClassesP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "classes")
}

func SRListDataTypes(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "datatypes")
}

func SRListDataTypesP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "datatypes")
}

func SRListDescriptors(ctx context.Context, a *AuthenticationConfig, p *SRListDescriptorsParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "descriptors")
}

func SRListDescriptorsP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "descriptors")
}

func SRListFieldGroups(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "fieldgroups")
}

func SRListFieldGroupsP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "fieldgroups")
}

func SRListSchemas(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "schemas")
}

func SRListSchemasP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "schemas")
}

func SRListUnions(ctx context.Context, a *AuthenticationConfig, p *SRListParams) (*http.Response, error) {
	return srList(ctx, a, p.Params(), "unions")
}

func SRListUnionsP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srList(ctx, a, p, "unions")
}

type SRGetBaseParams struct {
	ID     string
	Global bool
}

func (p *SRGetBaseParams) Params() util.Params {
	var cid string
	if p.Global {
		cid = "global"
	} else {
		cid = "tenant"
	}
	return util.NewParams(
		"-id", p.ID,
		"-cid", cid,
	)
}

type SRGetParams struct {
	SRGetBaseParams
	SRFormat
}

func (p *SRGetParams) Params() util.Params {
	result := p.SRGetBaseParams.Params()
	result.Add("-accept", p.SRFormat.Accept())
	return result
}

type SRGetDescriptorParams struct {
	SRGetBaseParams
	SRDescriptorFormat
}

func (p *SRGetDescriptorParams) Params() util.Params {
	result := p.SRGetBaseParams.Params()
	result.Add("-accept", p.SRDescriptorFormat.Accept())
	return result
}

func srGet(ctx context.Context, a *AuthenticationConfig, p util.Params, res string) (*http.Response, error) {
	header := map[string]string{"Accept": p.Get("-accept")}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/%s/%s/%s", p.GetForPath("-cid"), res, p.GetForPath("-id"))
}

func SRGetClass(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "classes")
}

func SRGetClassP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "classes")
}

func SRGetFieldGroup(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "fieldgroups")
}

func SRGetFieldGroupP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "fieldgroups")
}

func SRGetDataType(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "datatypes")
}

func SRGetDataTypeP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "datatypes")
}

func SRGetDescriptor(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "descriptors")
}

func SRGetDescriptorP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "descriptors")
}

func SRGetSchema(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "schemas")
}

func SRGetSchemaP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "schemas")
}

func SRGetUnion(ctx context.Context, a *AuthenticationConfig, p *SRGetParams) (*http.Response, error) {
	return srGet(ctx, a, p.Params(), "unions")
}

func SRGetUnionP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return srGet(ctx, a, p, "unions")
}

func SRGetSample(ctx context.Context, a *AuthenticationConfig, schemaID string) (*http.Response, error) {
	header := map[string]string{"Accept": "application/vnd.adobe.xed+json; version=1"}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/rpc/sampledata/%s", url.PathEscape(schemaID))
}

func SRGetAuditLog(ctx context.Context, a *AuthenticationConfig, schemaID string) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/rpc/auditlog/%s", url.PathEscape(schemaID))
}

type SRGetGlobalParams struct {
	SRFormat
	ID string
}

func (p *SRGetGlobalParams) Params() util.Params {
	return util.NewParams(
		"id", p.ID,
		"accept", p.Accept(),
	)
}

func SRGetBehavior(ctx context.Context, a *AuthenticationConfig, p *SRGetGlobalParams) (*http.Response, error) {
	return SRGetBehaviorP(ctx, a, p.Params())
}

func SRGetBehaviorP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	header := map[string]string{"Accept": p.Get("accept")}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/global/behaviors/%s", p.GetForPath("id"))
}

func SRExport(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return SRExportP(ctx, a, util.NewParams("id", id))
}

func SRExportP(ctx context.Context, a *AuthenticationConfig, p util.Params) (*http.Response, error) {
	header := map[string]string{"Accept": "application/vnd.adobe.xed-full+json; version=1"}
	return a.GetRequestHRaw(ctx, header, "https://platform.adobe.io/data/foundation/schemaregistry/rpc/export/%s", p.GetForPath("id"))
}

func SRImport(ctx context.Context, a *AuthenticationConfig, resource []byte) (*http.Response, error) {
	header := map[string]string{"Content-Type": "application/vnd.adobe.xed-full+json; version=1"}
	return a.PostRequestRaw(ctx, header, resource,
		"https://platform.adobe.io/data/foundation/schemaregistry/rpc/import")
}

func SRImportStream(ctx context.Context, a *AuthenticationConfig, r io.Reader) (*http.Response, error) {
	header := map[string]string{"Content-Type": "application/vnd.adobe.xed-full+json; version=1"}
	return a.PostRequestStream(ctx, header, r,
		"https://platform.adobe.io/data/foundation/schemaregistry/rpc/import")
}

func srDelete(ctx context.Context, a *AuthenticationConfig, resource, id string) (*http.Response, error) {
	url.PathEscape(resource)
	return a.DeleteRequestRaw(ctx,
		"https://platform.adobe.io/data/foundation/schemaregistry/tenant/%s/%s",
		url.PathEscape(resource),
		url.PathEscape(id))
}

func SRDeleteClass(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return srDelete(ctx, a, "classes", id)
}

func SRDeleteDataType(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return srDelete(ctx, a, "datatypes", id)
}

func SRDeleteDescriptor(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return srDelete(ctx, a, "descriptors", id)
}

func SRDeleteFieldGroup(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return srDelete(ctx, a, "fieldgroups", id)
}

func SRDeleteSchema(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return srDelete(ctx, a, "schemas", id)
}
