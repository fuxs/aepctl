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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/fuxs/aepctl/util"
	"github.com/rs/zerolog/log"
)

// BearerToken contains a token and an expiration date
type BearerToken struct {
	Token   string
	Expires time.Time
}

// ValidIn in checks if the token is still valid for passed duration
func (t *BearerToken) ValidIn(d time.Duration) bool {
	return time.Now().Add(d).Before(t.Expires)
}

// LocalTime returns the expiry date as string in the local time zone
func (t *BearerToken) LocalTime() string {
	return t.Expires.Local().Format(time.RFC822)
}

// AuthenticationConfig contains the configuraion for getting the bearer token
type AuthenticationConfig struct {
	Cache            bool
	DryRun           bool
	Server           string
	Organization     string
	TechnicalAccount string
	Audience         string
	ClientID         string
	ClientSecret     string
	Key              string
	Sandbox          string
	LoadToken        func() (*BearerToken, error)
	SaveToken        func(token *BearerToken) error
}

// UpdateHeader adds the authentication headers to the passed http request
func (o *AuthenticationConfig) UpdateHeader(req *http.Request) error {
	token, err := o.GetToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("x-api-key", o.ClientID)
	req.Header.Set("x-gw-ims-org-id", o.Organization)
	req.Header.Set("x-sandbox-name", o.Sandbox)
	return nil
}

func handleErrorResponse(res *http.Response) error {
	var sb strings.Builder
	sb.WriteString("Error (")
	sb.WriteString(res.Status)
	sb.WriteString("): ")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New(sb.String())
	}
	if q, err := util.UnmarshallQuery(body); err == nil {
		if s := q.Str("error_description"); s != "" {
			sb.WriteString(s)
			return errors.New(sb.String())
		}
		if s := q.Str("error"); s != "" {
			sb.WriteString(s)
			return errors.New(sb.String())
		}
	}
	sb.Write(body)
	return errors.New(sb.String())
}

// GetToken uses JWT to get a bearer token
func (o *AuthenticationConfig) GetToken() (*BearerToken, error) {
	if o.Cache && o.LoadToken != nil {
		if token, _ := o.LoadToken(); token != nil {
			if token.ValidIn(time.Minute) {
				return token, nil
			}
		}
	}
	//
	// build audience string
	audience := o.Audience
	if audience == "" {
		audience = "https://ims-na1.adobelogin.com/c/" + o.ClientID
	}
	//
	// create claim
	claim := &util.Claim{
		Iss: o.Organization,     // organization id
		Sub: o.TechnicalAccount, // technical account id
		Aud: audience,           // uses client id
	}
	//
	// load private key
	pKey, err := util.LoadPrivateKeyPEM(o.Key)
	if err != nil {
		return nil, err
	}
	//
	// sign claim with key
	token, err := claim.JWT(pKey)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"client_id":     {o.ClientID},
		"client_secret": {o.ClientSecret},
		"jwt_token":     {token},
	}
	server := o.Server
	if server == "" {
		server = "https://ims-na1.adobelogin.com/ims/exchange/jwt/"
	}
	res, err := http.PostForm(server, values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, handleErrorResponse(res)
	}
	var data []byte
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}{}

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	expiresIn := time.Millisecond * time.Duration(response.ExpiresIn)
	expires := time.Now().Add(expiresIn)

	result := &BearerToken{
		Token:   response.AccessToken,
		Expires: expires,
	}
	if o.Cache && o.SaveToken != nil {
		_ = o.SaveToken(result) // ignore error
	}

	return result, nil
}

// PostJSONRequest serializes the passed object to JSON and sends a httep post
// request to the passed url
func (o *AuthenticationConfig) PostJSONRequest(ctx context.Context, obj interface{}, url string, a ...interface{}) (interface{}, error) {
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return o.FullRequest(ctx, "POST", map[string]string{"Content-Type": "application/json"}, bytes.NewBuffer(body), url, a...)
}

// PostRequest sends a http post request to the passed url
func (o *AuthenticationConfig) PostRequest(ctx context.Context, header map[string]string, data []byte, url string, a ...interface{}) (interface{}, error) {
	return o.FullRequest(ctx, "POST", header, bytes.NewBuffer(data), url, a...)
}

// GetRequest sends a http get request to the passed url
func (o *AuthenticationConfig) GetRequest(ctx context.Context, url string, a ...interface{}) (interface{}, error) {
	return o.Request(ctx, "GET", url, a...)
}

// DeleteRequest sends a http delete request to the passed url
func (o *AuthenticationConfig) DeleteRequest(ctx context.Context, url string, a ...interface{}) (interface{}, error) {
	return o.Request(ctx, "DELETE", url, a...)
}

// PatchRequest sends a http patch request to the passed url
func (o *AuthenticationConfig) PatchRequest(ctx context.Context, header map[string]string, data []byte, url string, a ...interface{}) (interface{}, error) {
	return o.FullRequest(ctx, "PATCH", header, bytes.NewBuffer(data), url, a...)
}

// Request sends a http request with the passed verb to the passed url
func (o *AuthenticationConfig) Request(ctx context.Context, verb, url string, a ...interface{}) (interface{}, error) {
	return o.FullRequest(ctx, verb, nil, nil, url, a...)
}

// FullRequest sends a http request with the passed verb to the passed url
func (o *AuthenticationConfig) FullRequest(ctx context.Context, verb string, header map[string]string, body io.Reader, url string, a ...interface{}) (interface{}, error) {
	req, err := http.NewRequest(verb, fmt.Sprintf(url, a...), body)
	if err != nil {
		return nil, err
	}

	if err = o.UpdateHeader(req); err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	/*if verb == "POST" || verb == "PUT" || verb == "PATCH" {
		req.Header.Add("Content-Type", "application/json")
	}*/

	req = req.WithContext(ctx)

	if log.Debug().Enabled() {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		log.Debug().Str("Request", string(requestDump)).Msg("Dumping http request")
	}
	if o.DryRun {
		log.Debug().Msg("Dry-run done")
		//debug.PrintStack()
		return nil, nil
	}

	c := http.Client{
		Timeout: time.Minute,
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, handleErrorResponse(res)
	}
	var data []byte
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if log.Debug().Enabled() {
		log.Debug().Int("Code", res.StatusCode).Str("Status", res.Status).Str("Body", string(data)).Msg("Dumping http response")
	}

	if len(data) == 0 {
		return nil, nil
	}
	var obj map[string]interface{}
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return obj, nil

}
