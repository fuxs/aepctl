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
	"errors"
	"net/http"

	"github.com/fuxs/aepctl/util"
)

// DAGetFiles returns a list of files for the passed batchId
func DAGetFiles(ctx context.Context, p *AuthenticationConfig, batchId, start, limit string) (*http.Response, error) {
	if batchId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/batches/%s/files%s", batchId, util.Par("start", start, "limit", limit))
}

// DAGetFile returns a list of files for the passed fileId
func DAGetFile(ctx context.Context, p *AuthenticationConfig, fileId, start, limit string) (*http.Response, error) {
	if fileId == "" {
		return nil, errors.New("parameter fileId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, util.Par("start", start, "limit", limit))
}

func DADownload(ctx context.Context, p *AuthenticationConfig, fileId, path string) (*http.Response, error) {
	if fileId == "" {
		return nil, errors.New("parameter fileId is empty")
	}
	if path == "" {
		return nil, errors.New("url is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, util.Par("path", path))
}
