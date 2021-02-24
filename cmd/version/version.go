/*
Package version contains the version comman

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
package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Commit will be updated with git commit info
	Commit string = "unknown"
	// BuildTime will be updated with the build time
	BuildTime string = "unknown"
)

// NewCommand creates an initialized command object
func NewCommand(v string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "version",
		Short:                 "Returns version information",
		Long:                  "Returns version information",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("aepctl", v)
			fmt.Println("git commit:", Commit)
			fmt.Println("build time:", BuildTime)
			fmt.Println("go version:", runtime.Version())
			fmt.Println("platform  :", runtime.GOOS+"/"+runtime.GOARCH)
		},
	}
	return cmd
}
