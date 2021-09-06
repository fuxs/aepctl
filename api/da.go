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
)

type DAOptions struct {
	ID    string
	Start int
	Limit int
}

func (d *DAOptions) Request() *Request {
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
	req := NewRequest(
		"start", start,
		"limit", limit,
	)
	req.SetValue("id", d.ID)
	return req
}

// DAGetFilesP returns a list of files for the passed batchId
func DAGetFiles(ctx context.Context, p *AuthenticationConfig, options *DAOptions) (*http.Response, error) {
	return DAGetFilesP(ctx, p, options.Request())
}

// DAGetFilesP returns a list of files for the passed batchId
func DAGetFilesP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	batchId := params.GetValuePath("id")
	if batchId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/batches/%s/files%s", batchId, params.EncodedQuery())
}

// DAGetFile returns a list of files for the passed fileId
func DAGetFile(ctx context.Context, p *AuthenticationConfig, options *DAOptions) (*http.Response, error) {
	return DAGetFileP(ctx, p, options.Request())
}

// DAGetFileP returns a list of files for the passed fileId
func DAGetFileP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	fileId := params.GetValuePath("id")
	if fileId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, params.EncodedQuery())
}

// DADownload downloads a file fromt the given url
func DADownload(ctx context.Context, p *AuthenticationConfig, fileId, path string) (*http.Response, error) {
	req := NewRequest("path", path)
	req.SetValue("id", fileId)
	return DADownloadP(ctx, p, req)
}

// DADownload downloads a file fromt the given url
func DADownloadP(ctx context.Context, p *AuthenticationConfig, params *Request) (*http.Response, error) {
	fileId := params.GetValuePath("id")
	if fileId == "" {
		return nil, errors.New("parameter batchId is empty")
	}
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/export/files/%s%s", fileId, params.EncodedQuery())
}
