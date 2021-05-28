/*
Package util util consists of general utility functions and structures.

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
package util

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ConsoleWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	a := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(a) != 2 {
		return 0, errors.New("could not extract width")
	}
	i, err := strconv.Atoi(a[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}
