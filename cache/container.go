/*
Package cache consists of all caching relted functions and data structures.

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
package cache

import (
	"time"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// NewContainerCache creates a new container cache
func NewContainerCache(auth *api.AuthenticationConfig, pp util.PathProvider) *MapFileCache {
	t := NewTransformMap("_embedded", "https://ns.adobe.com/experience/xcore/container").
		K("_instance", "parentName").
		V("instanceId")
	return NewMapFileCache(NewContainerCall(auth), t, time.Hour*24, util.NewLazyPath(pp, "container.json"))
}

// AutoContainer resolves autmatically the container id
type AutoContainer struct {
	ContainerID string
	Auth        *api.AuthenticationConfig
	cc          *MapFileCache
}

func NewAutoContainer(auth *api.AuthenticationConfig, pp util.PathProvider) *AutoContainer {
	return &AutoContainer{Auth: auth, cc: NewContainerCache(auth, pp)}
}

// AddContainerFlag adds the required container flags to the passed command
func (a *AutoContainer) AddContainerFlag(cmd *cobra.Command) error {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&a.ContainerID, "container", "c", "", "Container where the decision rules are located")
	if err := cmd.RegisterFlagCompletionFunc("container", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return a.cc.Values(), cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		return err
	}
	return nil
}

// Get retrieves the container id
func (a *AutoContainer) Get() (string, error) {
	if a.ContainerID == "" {
		id, err := a.cc.LookupE(a.Auth.Sandbox)
		if err != nil {
			return "", err
		}
		a.ContainerID = id
	}
	return a.ContainerID, nil
}
