/*
Package extern is the root package for external commands.

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
package extern

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/cmd/helper"
	"github.com/spf13/cobra"
)

func NewPSQLCommand(conf *helper.Configuration) *cobra.Command {
	var (
		command string
		print   bool
	)
	cmd := &cobra.Command{
		Use:                   "psql",
		Short:                 "Run the psql client (Query Service)",
		Long:                  "long",
		Example:               "example",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			helper.CheckErrs(conf.Validate(cmd))
			q, err := api.NewQuery(api.HandleStatusCode(api.QSGetConnection(context.Background(), conf.Authentication)))
			helper.CheckErr(err)
			var sb strings.Builder
			if print {
				sb.WriteString(command)
				sb.WriteString(` "`)
			}
			sb.WriteString(`sslmode=require host=`)
			sb.WriteString(q.Str("host"))
			sb.WriteString(` port=`)
			sb.WriteString(q.Str("port"))
			sb.WriteString(` dbname=`)
			sb.WriteString(q.Str("dbName"))
			sb.WriteString(` user=`)
			sb.WriteString(q.Str("username"))
			sb.WriteString(` password=`)
			sb.WriteString(q.Str("token"))
			if print {
				sb.WriteRune('"')
				fmt.Println(sb.String())
			} else {
				c := exec.Command(command, sb.String())
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				helper.CheckErr(c.Run())
			}
		},
	}
	conf.AddAuthenticationFlags(cmd)
	flags := cmd.Flags()
	flags.StringVar(&command, "command", "psql", "path to psql command")
	flags.BoolVar(&print, "print", false, "don't exec but print command")
	return cmd
}
