/*
Package catalog consists of catalog functions.

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
)

// BatchesOptions contains all options for batch requests
type BatchesOptions struct {
	Limit         int
	TimeFormat    string
	CreatedAfter  string
	CreatedBefore string
	Dataset       string
	EndAfter      string
	EndBefore     string
	Name          string
	OrderBy       string
	StartAfter    string
	StartBefore   string
}

// ToURLPar converts the options to a URL parameter query
func (b *BatchesOptions) Request() (*Request, error) {
	var (
		limit         string
		createdAfter  string
		createdBefore string
		endAfter      string
		endBefore     string
		startAfter    string
		startBefore   string
	)
	tf := b.TimeFormat
	if tf == "" {
		tf = time.RFC3339
	}
	if b.Limit > 0 {
		limit = strconv.FormatInt(int64(b.Limit), 10)
	}
	if b.CreatedAfter != "" {
		t, err := time.Parse(tf, b.CreatedAfter)
		if err != nil {
			return nil, err
		}
		createdAfter = strconv.FormatInt(int64(t.Unix())*1000, 10)
	}
	if b.CreatedBefore != "" {
		t, err := time.Parse(tf, b.CreatedBefore)
		if err != nil {
			return nil, err
		}
		createdBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.EndAfter != "" {
		t, err := time.Parse(tf, b.EndAfter)
		if err != nil {
			return nil, err
		}
		endAfter = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.EndBefore != "" {
		t, err := time.Parse(tf, b.EndBefore)
		if err != nil {
			return nil, err
		}
		endBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.StartAfter != "" {
		t, err := time.Parse(tf, b.StartAfter)
		if err != nil {
			return nil, err
		}
		startAfter = strconv.FormatInt(t.Unix()*1000, 10)
	}
	if b.StartBefore != "" {
		t, err := time.Parse(tf, b.StartAfter)
		if err != nil {
			return nil, err
		}
		startBefore = strconv.FormatInt(t.Unix()*1000, 10)
	}
	return NewRequest(
		"limit", limit,
		"createdAfter", createdAfter,
		"createdBefore", createdBefore,
		"dataSet", b.Dataset,
		"endAfter", endAfter,
		"endBefore", endBefore,
		"name", b.Name,
		"orderBy", b.OrderBy,
		"startAfter", startAfter,
		"startBefore", startBefore,
	), nil
}

// CatalogGetBatches returns a list of batches
func CatalogGetBatches(ctx context.Context, p *AuthenticationConfig, options *BatchesOptions) (*http.Response, error) {
	params, err := options.Request()
	if err != nil {
		return nil, err
	}
	return CatalogGetBatchesP(ctx, p, params)
}

// CatalogGetBatchesP returns a list of batches
func CatalogGetBatchesP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/catalog/batches%s", params.EncodedQuery())
}

// CatalogGetDatasets returns a list of batches
func CatalogGetDatasets(ctx context.Context, p *AuthenticationConfig, options *BatchesOptions) (*http.Response, error) {
	params, err := options.Request()
	if err != nil {
		return nil, err
	}
	return CatalogGetDatasetsP(ctx, p, params)
}

// CatalogGetDatasets returns a list of batches
func CatalogGetDatasetsP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/catalog/datasets%s", params.EncodedQuery())
}
