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

	"github.com/fuxs/aepctl/util"
)

type FlowGetConnectionsParams struct {
	Property          string
	Limit             int
	OrderBy           string
	ContinuationToken string
	Count             bool
}

func (p *FlowGetConnectionsParams) Params() util.Params {
	var limit, count string
	if p.Limit > 0 {
		limit = strconv.FormatInt(int64(p.Limit), 10)
	}
	if p.Count {
		count = "true"
	}
	return util.NewParams(
		"property", p.Property,
		"limit", limit,
		"oderby", p.OrderBy,
		"continuationToken", p.ContinuationToken,
		"count", count,
	)
}

// DAGetFile returns a list of files for the passed fileId
func FlowGetConnections(ctx context.Context, auth *AuthenticationConfig, params util.Params) (*http.Response, error) {
	return auth.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/foundation/flowservice/connections%s",
		params.Encode())
}
