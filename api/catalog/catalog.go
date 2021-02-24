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

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// BatchesOptions contains all options for batch requests
type BatchesOptions struct {
	Limit        string
	CreatedAfter string
}

// ToURLPar converts the options to a URL parameter query
func (b *BatchesOptions) ToURLPar() string {
	return util.Par("limit", b.Limit, "createdAfter", b.CreatedAfter)
}

// GetBatches returns a list of batches
func GetBatches(ctx context.Context, p *api.AuthenticationConfig, o *BatchesOptions) (util.JSONResponse, error) {
	res, err := p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/catalog/batches%s", o.ToURLPar())
	if err != nil {
		return nil, err
	}
	return util.NewJSONMapIterator(res.Body)

	/*if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	i, err := util.NewJSONMapIterator(res.Body)
	if err != nil {
		return nil, err
	}
	for i.More() {
		n, err := i.Next()
		if err != nil {
			return nil, err
		}
		q := util.NewQuery(n)
		q.RangeAttributes(func(s string, q *util.Query) {
			fmt.Print(s)
			fmt.Println(" " + q.Str("createdUser"))
		})

	}
	return p.GetRequest(ctx, "https://platform.adobe.io/data/foundation/catalog/batches%s", o.ToURLPar())*/
}
