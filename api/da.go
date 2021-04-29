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
	"strconv"

	"github.com/fuxs/aepctl/util"
)

type DAOptions struct {
	ID    string
	Start int
	Limit int
}

func (d *DAOptions) Params() util.Params {
	var (
		limit string
		start string
	)
	if d.Start > 0 {
		start = strconv.FormatInt(int64(d.Start), 10)
	}
	if d.Limit > 0 {
		limit = strconv.FormatInt(int64(d.Limit), 10)
	}
	return util.NewParams(
		"-id", d.ID,
		"start", start,
		"limit", limit,
	)
}

// DAGetFilesP returns a list of files for the passed batchId
func DAGetFiles(ctx context.Context, p *AuthenticationConfig, options *DAOptions) (*http.Response, error) {
	return DAGetFilesP(ctx, p, options.Params())
}

// DAGetFilesP returns a list of files for the passed batchId
func DAGetFilesP(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	batchId := params.GetForPath("-id")
	if batchId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/batches/%s/files%s", batchId, params.EncodeWithout("-id"))
}

// DAGetFile returns a list of files for the passed fileId
func DAGetFile(ctx context.Context, p *AuthenticationConfig, options *DAOptions) (*http.Response, error) {
	return DAGetFileP(ctx, p, options.Params())
}

// DAGetFileP returns a list of files for the passed fileId
func DAGetFileP(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	fileId := params.GetForPath("-id")
	if fileId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, params.EncodeWithout("-id"))
}

// DADownload downloads a file fromt the given url
func DADownload(ctx context.Context, p *AuthenticationConfig, fileId, path string) (*http.Response, error) {
	return DADownloadP(ctx, p, util.NewParams("-id", fileId, "path", path))
}

// DADownload downloads a file fromt the given url
func DADownloadP(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	fileId := params.GetForPath("-id")
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, params.EncodeWithout("-id"))
}
