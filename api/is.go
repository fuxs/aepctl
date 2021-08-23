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
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/fuxs/aepctl/util"
)

func ISListNamespaces(ctx context.Context, a *AuthenticationConfig) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/core/idnamespace/identities")
}

func ISListNamespacesIMSOrg(ctx context.Context, a *AuthenticationConfig, imsOrg string) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/core/idnamespace/orgs/%s/identities", url.PathEscape(imsOrg))
}

func ISGetNamespace(ctx context.Context, a *AuthenticationConfig, id string) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/core/idnamespace/identities/%s", url.PathEscape(id))
}

func ISGetNamespaceIMSOrg(ctx context.Context, a *AuthenticationConfig, imsOrg, id string) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/core/idnamespace/orgs/%s/identities/%s", url.PathEscape(imsOrg), url.PathEscape(id))
}

func ISCreateNamespace(ctx context.Context, a *AuthenticationConfig, body []byte) (*http.Response, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return a.PostRequestRaw(ctx, header, body, "https://platform.adobe.io/data/core/idnamespace/identities")
}

func ISUpdateNamespace(ctx context.Context, a *AuthenticationConfig, id string, body []byte) (*http.Response, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return a.PutRequestRaw(ctx, header, body, "https://platform.adobe.io/data/core/idnamespace/identities/%s", url.PathEscape(id))
}

type ISParams struct {
	Namespace   string
	NamespaceID string
	ID          string
	Region      string
}

func (p *ISParams) Validate() error {
	if p.Namespace == "" && p.NamespaceID == "" {
		return errors.New("no namespace id or namespace set")
	}
	if p.Namespace != "" && p.NamespaceID != "" {
		return errors.New("namespace id and namespace code are used together")
	}
	if len(p.ID) == 0 {
		return errors.New("no ID set")
	}
	return nil
}

func (p *ISParams) Request() *Request {
	idName := "xid"
	if p.Namespace != "" || p.NamespaceID != "" {
		idName = "id"
	}

	req := NewRequest(
		idName, p.ID,
		"namespace", p.Namespace,
		"nsid", p.NamespaceID,
	)
	req.SetValue("region", p.Region)
	return req
}

func ISGetXIDR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	return a.GetRequestRaw(ctx,
		"https://platform-%s.adobe.io/data/core/identity/identity%s",
		req.GetValueV("region", "va7"),
		req.EncodedQuery())
}

// ISGetXID returns the XID for the given namespace and id combination
func ISGetXID(ctx context.Context, a *AuthenticationConfig, p *ISParams) (*http.Response, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return ISGetXIDR(ctx, a, p.Request())
}

type ISClusterParams struct {
	ISParams
	GraphType string
}

func (p *ISClusterParams) Request() *Request {
	result := p.ISParams.Request()
	result.Query.Add("graph-type", p.GraphType)
	return result
}

func (p *ISClusterParams) Params() util.Params {
	idName := "xid"
	if p.Namespace != "" || p.NamespaceID != "" {
		idName = "id"
	}
	return util.NewParams(
		idName, p.ID,
		"namespace", p.Namespace,
		"nsid", p.NamespaceID,
	)
}

func (p *ISClusterParams) Accept() string {
	return "application/vnd.adobe.identity+json;version=1.2"
}

func ISGetClusterR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	return a.GetRequestRaw(ctx,
		"https://platform-%s.adobe.io/data/core/identity/cluster/members%s",
		req.GetValueV("region", "va7"),
		req.EncodedQuery())
}

func ISGetCluster(ctx context.Context, a *AuthenticationConfig, p *ISClusterParams) (*http.Response, error) {
	return ISGetClusterR(ctx, a, p.Request())
}

func ISGetHistoryR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	return a.GetRequestRaw(ctx,
		"https://platform-%s.adobe.io/data/core/identity/cluster/history%s",
		req.GetValueV("region", "va7"),
		req.EncodedQuery())
}

func ISGetHistory(ctx context.Context, a *AuthenticationConfig, p *ISClusterParams) (*http.Response, error) {
	return ISGetHistoryR(ctx, a, p.Request())
}

type ISClustersParams struct {
	Namespaces   []string
	NamesapceIDs []string
	IDs          []string
	GraphType    string
	Region       string
}

func (p *ISClustersParams) Body() ([]byte, error) {
	var sb bytes.Buffer
	l := len(p.IDs)
	lns := len(p.Namespaces)
	lid := len(p.NamesapceIDs)
	if lns > 0 && lid > 0 {
		return nil, errors.New("input error: namespace ids and namespace codes are used together")
	}
	if lns+lid > l {
		return nil, errors.New("input error: more namespaces than ids")
	}
	if lns == 0 && lid == 0 {
		sb.WriteString(`{"xids": [`)
		for i, id := range p.IDs {
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteRune('"')
			sb.WriteString(id)
			sb.WriteRune('"')
		}
	} else {
		sb.WriteString(`{"compositeXids": [`)
		var ns string
		if lid > lns {
			for i, id := range p.IDs {
				if i > 0 {
					sb.WriteRune(',')
				}
				sb.WriteString(`{"nsid": `)
				if i < lid {
					ns = p.NamesapceIDs[i]
				}
				sb.WriteString(ns)
				sb.WriteString(`,"id":"`)
				sb.WriteString(id)
				sb.WriteString(`"}`)
			}
		} else {
			for i, id := range p.IDs {
				if i > 0 {
					sb.WriteRune(',')
				}
				sb.WriteString(`{"namespace": "`)
				if i < lns {
					ns = p.Namespaces[i]
				}
				sb.WriteString(ns)
				sb.WriteString(`","id":"`)
				sb.WriteString(id)
				sb.WriteString(`"}`)
			}
		}
	}
	sb.WriteRune(']')
	if p.GraphType != "" {
		sb.WriteString(`,"graph-type": "`)
		sb.WriteString(p.GraphType)
		sb.WriteRune('"')
	}
	sb.WriteRune('}')
	return sb.Bytes(), nil
}

func (p *ISClustersParams) Request() (*Request, error) {
	body, err := p.Body()
	if err != nil {
		return nil, err
	}
	return NewRequestBody(body), nil
}

func ISGetClustersR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	req.ContentType("application/json")
	return a.PostRequestRaw(ctx,
		req.Header(),
		req.Body,
		"https://platform-%s.adobe.io/data/core/identity/clusters/members",
		req.GetValueV("region", "va7"))
}

func ISGetClusters(ctx context.Context, a *AuthenticationConfig, p *ISClustersParams) (*http.Response, error) {
	req, err := p.Request()
	if err != nil {
		return nil, err
	}
	return ISGetClustersR(ctx, a, req)
}

func ISGetHistoriesR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	req.ContentType("application/json")
	return a.PostRequestRaw(ctx,
		req.Header(),
		req.Body,
		"https://platform-%s.adobe.io/data/core/identity/clusters/history",
		req.GetValueV("region", "va7"))
}

func ISGetHistories(ctx context.Context, a *AuthenticationConfig, p *ISClustersParams) (*http.Response, error) {
	req, err := p.Request()
	if err != nil {
		return nil, err
	}
	return ISGetHistoriesR(ctx, a, req)
}

type ISGetMappingParams struct {
	ISParams
	TargetNS string
}

func ISGetMappingR(ctx context.Context, a *AuthenticationConfig, req *Request) (*http.Response, error) {
	return a.GetRequestRaw(ctx,
		"https://platform-%s.adobe.io/data/core/identity/mapping%s",
		req.GetValueV("region", "va7"),
		req.EncodedQuery())
}

func ISGetMapping(ctx context.Context, a *AuthenticationConfig, p *ISGetMappingParams) (*http.Response, error) {
	return ISGetMappingR(ctx, a, p.Request())
}

func ISGetMappings(ctx context.Context, a *AuthenticationConfig, p *ISClustersParams) (*http.Response, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	/*body, err := p.Body()
	if err != nil {
		return nil, err
	}*/
	body := `{
		"xids": [
		  {
			"xid": "A28eOco1-QqGQERvuJjKVoEc"
		  }
		],
		"targetNs": 4
	  }`

	return a.PostRequestRaw(ctx, header, []byte(body), "https://platform.adobe.io/data/core/identity/mapping")
}
