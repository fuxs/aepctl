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
	"time"

	"github.com/fuxs/aepctl/util"
)

type UPSEntitiesParams struct {
	Schema        string
	RelatedSchema string
	ID            string
	NS            string
	RelatedID     string
	RelatedNS     string
	Fields        string
	MP            string
	Start         string
	End           string
	Limit         int
	Order         string
	Property      string
	CA            bool
	TimeFormat    string
}

func (p *UPSEntitiesParams) Params() (util.Params, error) {
	var start, end, limit, ca string
	tf := p.TimeFormat
	if tf == "" {
		tf = time.RFC3339
	}
	if p.Start != "" {
		if util.IsNumeric(p.Start) {
			start = p.Start
		} else {
			t, err := time.Parse(tf, p.Start)
			if err != nil {
				return nil, err
			}
			start = strconv.FormatInt(t.Unix()*1000, 10)
		}
	}
	if p.End != "" {
		if util.IsNumeric(p.End) {
			start = p.End
		} else {
			t, err := time.Parse(tf, p.End)
			if err != nil {
				return nil, err
			}
			start = strconv.FormatInt(t.Unix()*1000, 10)
		}
	}
	if p.Limit > 0 {
		limit = strconv.FormatInt(int64(p.Limit), 10)
	}
	if p.CA {
		ca = "true"
	}
	return util.NewParams(
		"schema.name", p.Schema,
		"relatedSchema.name", p.RelatedSchema,
		"entityId", p.ID,
		"entityIdNS", p.NS,
		"relatedEntityId", p.RelatedID,
		"relatedEntityIdNS", p.RelatedNS,
		"fields", p.Fields,
		"mergePolicyId", p.MP,
		"startTime", start,
		"endTime", end,
		"limit", limit,
		"oderby", p.Order,
		"property", p.Property,
		"withCA", ca,
	), nil
}

func UPSGetEntities(ctx context.Context, auth *AuthenticationConfig, p *UPSEntitiesParams) (*http.Response, error) {
	params, err := p.Params()
	if err != nil {
		return nil, err
	}
	return auth.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/core/ups/access/entities%s",
		params.Encode(),
	)
}

func UPSGetEntitiesP(ctx context.Context, auth *AuthenticationConfig, p util.Params) (*http.Response, error) {
	return auth.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/core/ups/access/entities%s",
		p.Encode(),
	)
}
