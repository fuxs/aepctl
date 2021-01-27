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
	"github.com/fuxs/aepctl/util"
	"github.com/spf13/cobra"
)

// AutoContainer resolves autmatically the container id
type AutoContainer struct {
	ContainerID string
	Store       *util.KVCache
	Auth        *Authentication
}

// NewAutoContainer creates an initialized AutoContainer object
func NewAutoContainer(auth *Authentication) *AutoContainer {
	return &AutoContainer{
		Store: NewContainerCache(auth),
		Auth:  auth,
	}
}

// AddContainerFlag adds the required container flags to the passed command
func (q *AutoContainer) AddContainerFlag(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&q.ContainerID, "container", "c", "", "Container where the decision rules are located")
	if err := cmd.RegisterFlagCompletionFunc("container", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		containers, _ := q.Store.Values()
		return containers, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fatal("Error in AddContainerFlag", 1)
	}
}

// AutoFillContainer retrieves the container id
func (q *AutoContainer) AutoFillContainer() error {
	if q.ContainerID == "" {
		id, err := q.Store.GetValueE(q.Auth.Config.Sandbox)
		if err != nil {
			return err
		}
		q.ContainerID = id
	}
	return nil
}

// UniquePath generates a unique path based on the client id, sandbox name and container id
func (q *AutoContainer) UniquePath(path ...string) []string {
	if err := q.AutoFillContainer(); err != nil {
		return []string{}
	}
	auth := q.Auth.Config
	return append([]string{auth.ClientID, auth.Sandbox, q.ContainerID}, path...)
}
