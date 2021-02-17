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
	"errors"
	"strings"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// Configuration encapsulates the global settings
type Configuration struct {
	Root           *util.RootConfig
	Authentication *api.AuthenticationConfig
	Read           bool
	Write          bool
	AC             *AutoContainer
	AS             *util.KVCache
	TS             *util.KVCache
	PS             *util.KVCache
	OS             *util.KVCache
	FS             *util.KVCache
	CS             *util.KVCache
	RS             *util.KVCache
}

// NewConfiguration creates an initialized Authentication object
func NewConfiguration(gcfg *util.RootConfig) *Configuration {
	result := &Configuration{
		Root:           gcfg,
		Authentication: &api.AuthenticationConfig{},
	}
	cache := util.NewJSONCache(func() []string { return result.UniquePath("token.json") })
	o := result.Authentication
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

// Update loads the configuration and updates the command flags
func (a *Configuration) Update(cmd *cobra.Command) error {
	return a.Root.Configure(cmd)
}

// AddAuthenticationFlags adds all required flags for authentication
func (a *Configuration) AddAuthenticationFlags(cmd *cobra.Command) {
	o := a.Authentication
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
	flags.StringVar(&o.Key, "key", "private.key", "path to private key file")

	if err := cmd.RegisterFlagCompletionFunc("sandbox", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		sandboxes, _ := NewSandboxCache(a).GetList()
		return sandboxes, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fatal("Error in AddAuthenticationFlags", 1)
	}
}

// Validate updates the command flags and validates the final configuration
func (a *Configuration) Validate(cmd *cobra.Command) error {
	if err := a.Update(cmd); err != nil {
		return err
	}
	o := a.Authentication

	var clientID, clientSecret, techAccount, organization, key bool
	errCounter := 0
	if o.ClientID == "" {
		clientID = true
		errCounter++
	}
	if o.ClientSecret == "" {
		clientSecret = true
		errCounter++
	}
	if o.TechnicalAccount == "" {
		techAccount = true
		errCounter++
	}
	if o.Organization == "" {
		organization = true
		errCounter++
	}
	if o.Key == "" {
		key = true
		errCounter++
	}

	if errCounter > 0 {
		var (
			b     strings.Builder
			comma bool
		)
		if errCounter == 1 {
			b.WriteString("Missing authentication parameter ")
		} else {
			b.WriteString("Missing authentication parameters ")
		}
		if clientID {
			b.WriteString("Client ID (--client-id)")
			comma = true
		}
		if clientSecret {
			if comma {
				b.WriteString(", ")
			}
			comma = true
			b.WriteString("Client Secret (--client-secret)")
		}
		if techAccount {
			if comma {
				b.WriteString(", ")
			}
			comma = true
			b.WriteString("Technial Account ID (--tech-account)")
		}
		if organization {
			if comma {
				b.WriteString(", ")
			}
			comma = true
			b.WriteString("Organization ID (--organization)")
		}
		if key {
			if comma {
				b.WriteString(", ")
			}
			b.WriteString("Private Key File (--key)")
		}

		b.WriteString("\n")
		b.WriteString(util.Form(`
		Please provide all required flags or a configuration file (--config).

		Execute the following command to initialize aepctl:
		
		  aepctl configure`))

		return errors.New(b.String())
	}
	return nil
}

// NoDryRun creates a copy of the current AuthenticationConfig  and disables the dry-run falg
func (a *Configuration) NoDryRun() *api.AuthenticationConfig {
	cfg := *a.Authentication
	cfg.DryRun = false
	return &cfg
}

// UniquePath generates a unique path based on the client id
func (a *Configuration) UniquePath(path ...string) []string {
	return append([]string{a.Authentication.ClientID}, path...)
}

// UniqueSandboxPath generates a unique path based on the client id and sandbox name
func (a *Configuration) UniqueSandboxPath(path ...string) []string {
	return append([]string{a.Authentication.ClientID, a.Authentication.Sandbox}, path...)
}

// ReadCache returns the read-cache flag
func (a *Configuration) ReadCache() bool {
	return a.Read
}

// WriteCache returns the write-cache flag
func (a *Configuration) WriteCache() bool {
	return a.Write
}
