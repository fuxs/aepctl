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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fuxs/aepctl/api"
	"github.com/spf13/cobra"
)

// CheckErr prints a user friendly error message to stderr
func CheckErr(err error) {
	formatError(err, fatal)
}

func CheckErrInfo(err error) {
	formatError(err, info)
}

func CheckErrParams(params api.ParametersE) *api.Request {
	p, err := params.Request()
	if err != nil {
		formatError(err, fatal)
	}
	return p
}

// CheckErrEOF prints a user friendly error message to stderr but ignores EOF
func CheckErrEOF(err error) {
	if err == io.EOF {
		return
	}
	formatError(err, fatal)
}

func info(msg string, code int) {
	if len(msg) > 0 {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
}

func fatal(msg string, code int) {
	info(msg, code)
	os.Exit(code)
}

func formatError(err error, handler func(string, int)) {
	if err == nil {
		return
	}
	handler(err.Error(), 1)
}

// CheckErrs prints a user friendly error message to stderr
func CheckErrs(err ...error) {
	for _, e := range err {
		CheckErr(e)
	}
}

func PrintError(msg string, cmd *cobra.Command) {
	fmt.Fprintln(os.Stderr, msg)
	if err := cmd.Help(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(1)
}
