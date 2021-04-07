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
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/fuxs/aepctl/util"
)

func NewJSONIterator(res *http.Response, err error) (*util.JSONIterator, error) {
	if err != nil {
		return nil, err
	}
	return util.NewJSONIterator(res.Body), nil
}

func NewJSONFilterIterator(filter []string, res *http.Response, err error) (*util.JSONFilterIterator, error) {
	if err != nil {
		return nil, err
	}
	return util.NewJSONFilterIterator(filter, res.Body), nil
}

func NewQuery(res *http.Response, err error) (*util.Query, error) {
	i, err := NewJSONIterator(HandleStatusCode(res, err))
	if err != nil {
		return nil, err
	}
	return i.Query()
}

func HandleStatusCode(res *http.Response, err error) (*http.Response, error) {
	if err != nil || (res.StatusCode >= 200 && res.StatusCode < 300) {
		return res, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return nil, errors.New(string(data))
}
