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
	"strconv"
	"strings"
	"time"
)

// Difference substitues the array of strings b from a
func Difference(a, b []string) []string {
	amap := make(map[string]int)
	for _, s := range a {
		amap[s]++
	}
	for _, s := range b {
		delete(amap, s)
	}
	result := make([]string, 0, len(amap))
	for k := range amap {
		result = append(result, k)
	}
	return result
}

// RemoveQuotes removes embracing quotes, e.g. "hello" becomes hello and 'world' becomes world
func RemoveQuotes(s string) string {
	l := len(s) - 1
	if l >= 1 {
		if (s[0] == '\'' && s[l] == '\'') ||
			(s[0] == '"' && s[l] == '"') {
			s = s[1:l]
		}
	}
	return s
}

// AddDollar prepends a dollar sign to the string if it
func AddDollar(s string) string {
	if len(s) >= 1 {
		if s[0] != '$' {
			s = "$" + s
		}
	} else {
		s = "$"
	}
	return s
}

// Dashed returns the string "-" if passed string is empty.
func Dashed(str string) string {
	if str != "" {
		return str
	}
	return "-"
}

// GetString accepts a value and tries to return a string representation.
// If a conversion is not possible, then this function returns an error.
func GetString(value interface{}) (result string) {
	switch i := value.(type) {
	case nil:
		result = "-"
	case string:
		result = i
	case float64:
		result = strconv.FormatFloat(i, 'f', -1, 64)
	case float32:
		result = strconv.FormatFloat(float64(i), 'f', -1, 32)
	case int64:
		result = strconv.FormatInt(i, 10)
	case int32:
		result = strconv.FormatInt(int64(i), 10)
	case int:
		result = strconv.FormatInt(int64(i), 10)
	case int16:
		result = strconv.FormatInt(int64(i), 10)
	case int8:
		result = strconv.FormatInt(int64(i), 10)
	case uint64:
		result = strconv.FormatUint(i, 10)
	case uint32:
		result = strconv.FormatUint(uint64(i), 10)
	case uint:
		result = strconv.FormatUint(uint64(i), 10)
	case uint16:
		result = strconv.FormatUint(uint64(i), 10)
	case uint8:
		result = strconv.FormatUint(uint64(i), 10)
	case bool:
		result = strconv.FormatBool(i)
	default:
	}
	return
}

// LocalTimeStr converts time string in RFC3339 to local time in RFC822
func LocalTimeStr(str string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str
	}
	return t.Local().Format(time.RFC822)
}

// LocalTimeStrCustom converts time string in RFC3339 to local time in custom format
func LocalTimeStrCustom(str, fmt string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str
	}
	return t.Local().Format(fmt)
}

// GetInt returns the int value for the passed value or -1 for invalid input.
func GetInt(value interface{}) (result int) {
	switch i := value.(type) {
	case nil:
		result = 0
	case string:
		result, _ = strconv.Atoi(i)
	case float64:
		result = int(i)
	case float32:
		result = int(i)
	case int64:
		result = int(i)
	case int32:
		result = int(i)
	case int:
		result = i
	case int16:
		result = int(i)
	case int8:
		result = int(i)
	case uint64:
		result = int(i)
	case uint32:
		result = int(i)
	case uint:
		result = int(i)
	case uint16:
		result = int(i)
	case uint8:
		result = int(i)
	case bool:
		if i {
			result = 1
		} else {
			result = 0
		}
	default:
		result = -1
	}
	return
}

// Plural generates the English plural of the passed string
func Plural(str string) string {
	var result strings.Builder
	l := len(str)
	if l > 0 {
		if str[l-1] == 'y' {
			result.WriteString(str[:l-1])
			result.WriteString("ie")
		} else {
			result.WriteString(str)
		}
		result.WriteByte('s')
	}
	return result.String()
}

// Contains checks for the existence of the passed value in the array
func Contains(value string, ar []string) bool {
	for _, str := range ar {
		if value == str {
			return true
		}
	}
	return false
}

// ContainsS checks for the existence of the passed value in the array
func ContainsS(value string, ar []string) string {
	for _, str := range ar {
		if value == str {
			return "●"
		}
	}
	return "◯"
}
