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
	"io"
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

// DAGetFile returns a list of files for the passed fileId
func FlowGetConnections(ctx context.Context, auth *AuthenticationConfig, p *FlowGetConnectionsParams) (*http.Response, error) {
	var limit, count string
	if p.Limit > 0 {
		limit = strconv.FormatInt(int64(p.Limit), 10)
	}
	if p.Count {
		count = "true"
	}
	return auth.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/foundation/flowservice/connections%s",
		util.Par(
			"property", p.Property,
			"limit", limit,
			"oderby", p.OrderBy,
			"continuationToken", p.ContinuationToken,
			"count", count,
		))
}

func FlowGetNext(ctx context.Context, auth *AuthenticationConfig, url string) (*http.Response, error) {
	return auth.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/flowservice%s", url)
}

type Paged interface {
	First() (util.JSONResponse, error)
	Execute([]string, func(util.JSONResponse) error) error
}
type FlowPaged struct {
	ctx  context.Context
	auth *AuthenticationConfig
	p    *FlowGetConnectionsParams
}

func NewFlowPaged(ctx context.Context, auth *AuthenticationConfig, p *FlowGetConnectionsParams) *FlowPaged {
	return &FlowPaged{ctx: ctx, auth: auth, p: p}
}

func (fp *FlowPaged) First() (util.JSONResponse, error) {
	res, err := FlowGetConnections(fp.ctx, fp.auth, fp.p)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(data))
	}
	return util.NewJSONIterator(res.Body), nil
}

func (fp *FlowPaged) Execute(path []string, f func(util.JSONResponse) error) error {
	jf := util.NewJSONFinder()
	jf.Add(f, path...)
	next := true
	url := ""
	jf.Add(func(j util.JSONResponse) error {
		q, err := j.Query()
		if err != nil {
			return err
		}
		url = q.Str("next", "href")
		next = url != ""
		return nil
	}, "_links")

	// next and i are used by Run()
	for next {
		// anonymous function for innder defer i.Close commands
		func() error {
			next = false
			if url != "" {
				res, err := FlowGetNext(fp.ctx, fp.auth, url)
				if err != nil {
					return err
				}
				if res.StatusCode >= 300 {
					data, err := io.ReadAll(res.Body)
					if err != nil {
						return err
					}
					return errors.New(string(data))
				}
				i := util.NewJSONIterator(res.Body)
				defer i.Close()
				jf.SetIterator(i)
			} else {
				i, err := fp.First()
				if err != nil {
					return err
				}
				defer i.Close()
				jf.SetIterator(i)
			}
			if err := jf.Run(); err != nil {
				return err
			}
			return nil
		}()
	}
	return nil
}
