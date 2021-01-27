/*
Package helper consists of helping functions.

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
package helper

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// Authentication encapsulates the authentication configuration
type Authentication struct {
	Config *api.AuthenticationConfig
	Read   bool
	Write  bool
	AC     *AutoContainer
	AS     *util.KVCache
	TS     *util.KVCache
	PS     *util.KVCache
	OS     *util.KVCache
	FS     *util.KVCache
	CS     *util.KVCache
	RS     *util.KVCache
}

// NewAuthentication creates an initialized Authentication object
func NewAuthentication() *Authentication {
	result := &Authentication{
		Config: &api.AuthenticationConfig{},
	}
	ac := NewAutoContainer(result)
	result.AC = ac
	result.AS = NewActivityCache(ac)
	result.TS = NewTagCache(ac)

	result.PS = NewPlacementCache(ac)
	result.OS = NewOfferCache(ac)
	result.FS = NewFallbackCache(ac)
	result.CS = NewCollectionCache(ac)
	result.RS = NewRuleCache(ac)
	return result
}

// AddAuthenticationFlags adds all required flags for authentication
func (a *Authentication) AddAuthenticationFlags(cmd *cobra.Command) {
	o := a.Config
	flags := cmd.PersistentFlags()
	flags.BoolVar(&o.Cache, "cache", true, "stores the retrieved token in ~/.aepctl/token.json")
	flags.BoolVar(&a.Read, "read-cache", true, "stores the retrieved token in ~/.aepctl/token.json")
	flags.BoolVar(&a.Write, "write-cache", true, "stores the retrieved token in ~/.aepctl/token.json")

	flags.BoolVar(&o.DryRun, "dry-run", false, "builds the request but doesn't execute it.")
	flags.StringVar(&o.Server, "server", "https://ims-na1.adobelogin.com/ims/exchange/jwt/", "OAuth 2.0 server")
	flags.StringVar(&o.Organization, "organization", "", "organization")
	flags.StringVar(&o.TechnicalAccount, "tech-account", "", "technical account id")
	flags.StringVar(&o.Audience, "audience", "", "JWT audience")
	flags.StringVar(&o.ClientID, "client-id", "", "client id")
	flags.StringVar(&o.ClientSecret, "client-secret", "", "client secret")
	flags.StringVar(&o.Sandbox, "sandbox", "prod", "selects the sandbox (default is the name of the production sandbox: prod)")
	flags.StringVar(&o.Key, "key", "./extern/private.key", "path to private key file")
	cache := util.NewJSONCache(func() []string { return a.UniquePath("token.json") })
	o.LoadToken = func() (*api.BearerToken, error) {
		token := &api.BearerToken{}
		if err := cache.Load(token); err != nil {
			return nil, err
		}
		return token, nil
	}
	o.SaveToken = func(token *api.BearerToken) error {
		return cache.Save(token)
	}
	if err := cmd.RegisterFlagCompletionFunc("sandbox", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		sandboxes, _ := NewSandboxCache(a).GetList()
		return sandboxes, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fatal("Error in AddAuthenticationFlags", 1)
	}
}

// NoDryRun creates a copy of the current AuthenticationConfig  and disables the dry-run falg
func (a *Authentication) NoDryRun() *api.AuthenticationConfig {
	cfg := *a.Config
	cfg.DryRun = false
	return &cfg
}

// UniquePath generates a unique path based on the client id
func (a *Authentication) UniquePath(path ...string) []string {
	return append([]string{a.Config.ClientID}, path...)
}

// UniqueSandboxPath generates a unique path based on the client id and sandbox name
func (a *Authentication) UniqueSandboxPath(path ...string) []string {
	return append([]string{a.Config.ClientID, a.Config.Sandbox}, path...)
}

// ReadCache returns the read-cache flag
func (a *Authentication) ReadCache() bool {
	return a.Read
}

// WriteCache returns the write-cache flag
func (a *Authentication) WriteCache() bool {
	return a.Write
}
