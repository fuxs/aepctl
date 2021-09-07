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
)

// QSListQueriesParams defines the parameters for list queries
type QSListQueriesParams struct {
	PageParams
	ExcludeSoftDeleted bool
	ExcludeHidden      bool
}

func (p *QSListQueriesParams) Request() *Request {
	var esd, eh string
	if !p.ExcludeSoftDeleted {
		esd = "false"
	}
	if !p.ExcludeHidden {
		eh = "false"
	}
	req := p.PageParams.Request()
	req.AddQueries(
		"excludeSoftDeleted", esd,
		"excludeHidden", eh)
	return req
}

// QSListQueries calls the query servie to list queries
func QSListQueries(ctx context.Context, a *AuthenticationConfig, p *QSListQueriesParams) (*http.Response, error) {
	if p != nil {
		return QSListQueriesP(ctx, a, p.Request())
	}
	return QSListQueriesP(ctx, a, nil)
}

// QSListQueriesR calls the query servie to list queries
func QSListQueriesP(ctx context.Context, a *AuthenticationConfig, p *Request) (*http.Response, error) {
	if p == nil {
		return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/query/queries")
	}
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/query/queries%s", p.EncodedQuery())
}

// QSGetConnection retrieves connection parameters for the interactive interface
func QSGetConnection(ctx context.Context, a *AuthenticationConfig) (*http.Response, error) {
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/query/connection_parameters")
}

// QSListQueriesR calls the query servie to list queries
func QSListSchedulesP(ctx context.Context, a *AuthenticationConfig, p *Request) (*http.Response, error) {
	if p == nil {
		return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/query/schedules")
	}
	return a.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/query/schedules%s", p.EncodedQuery())
}

func QSCreateSchedule(ctx context.Context, a *AuthenticationConfig, body []byte) (*http.Response, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return a.PostRequestRaw(ctx, header, body, "https://platform.adobe.io/data/foundation/query/schedules")
}
