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
package catalog

import (
	"context"
	"net/http"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// BatchesOptions contains all options for batch requests
type BatchesOptions struct {
	Limit         string
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
func (b *BatchesOptions) ToURLPar() string {
	return util.Par(
		"limit", b.Limit,
		"createdAfter", b.CreatedAfter,
		"createdBefore", b.CreatedBefore,
		"dataSet", b.Dataset,
		"endAfter", b.EndAfter,
		"endBefore", b.EndBefore,
		"name", b.Name,
		"orderBy", b.OrderBy,
		"startAfter", b.StartAfter,
		"startBefore", b.StartBefore,
	)
}

// GetBatches returns a list of batches
func GetBatches(ctx context.Context, p *api.AuthenticationConfig, o *BatchesOptions) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/catalog/batches%s", o.ToURLPar())
}

// GetDatasets returns a list of batches
func GetDatasets(ctx context.Context, p *api.AuthenticationConfig, o *BatchesOptions) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/catalog/datasets%s", o.ToURLPar())
}
