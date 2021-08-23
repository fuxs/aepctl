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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/fuxs/aepctl/util"
)

func NewJSONIterator(res *http.Response, err error) (*util.JSONIterator, error) {
	if err != nil {
		return nil, err
	}
	return util.NewJSONIterator(util.NewJSONCursor(res.Body)), nil
}

func NewQuery(res *http.Response, err error) (*util.Query, error) {
	i, err := NewJSONIterator(HandleStatusCode(res, err))
	if err != nil {
		return nil, err
	}
	return i.Query()
}

func printPayload(res *http.Response, w io.Writer) error {
	mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if mt == "application/json" {
		if err := util.JSONPrintPretty(json.NewDecoder(res.Body), w); err != nil {
			return err
		}
	} else {
		_, err := io.Copy(w, res.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// HandleStatusCode checks for a previous error and the HTTP status code. If an
// error exists or the status code is good then the passed objects will be
// returned. Otherwise it will create an error object from the status
// information and the returned HTTP body.
func HandleStatusCode(res *http.Response, err error) (*http.Response, error) {
	if err != nil || (res.StatusCode >= 200 && res.StatusCode < 300) {
		return res, err
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, "http error with status code:", res.StatusCode)
	if err = printPayload(res, &sb); err != nil {
		return res, err
	}
	return res, errors.New(sb.String())
}

func DropResponse(res *http.Response, err error) error {
	_, err = HandleStatusCode(res, err)
	return err
}

func PrintResponse(res *http.Response, err error) error {
	// check for previous error
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Fprintln(os.Stderr, "http error with status code:", res.StatusCode)
	}
	if err = printPayload(res, os.Stdout); err != nil {
		return err
	}
	fmt.Println()
	return nil
}
