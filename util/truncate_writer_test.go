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
	"strings"
	"testing"
)

func TestNormal(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "123456789"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte("123456789")) = %v, %v, want %v, nil`, n, err, l)
	}
	result := sb.String()
	if result != str {
		t.Errorf(`sb.String() = %q,  want %q`, result, str)
	}

}

func TestNormalUnicode(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "12345678€"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != str {
		t.Errorf(`sb.String() = %q,  want %q`, result, str)
	}

}

func TestEmpty(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := ""
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte("")) = %v, %v, want %v, nil`, n, err, l)
	}
	result := sb.String()
	if result != str {
		t.Errorf(`sb.String() = %q,  want %q`, result, str)
	}
}

func TestMultiline(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "1234\n5678\n90ab"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != str {
		t.Errorf(`sb.String() = %q,  want %q`, result, str)
	}
}

func TestLong(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "1234567890ab"
	expect := "123456789…"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}

}

func TestMultiLong(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "1234567890abcdef\n1234567890abcde"
	expect := "123456789…\n123456789…"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}

}

func TestVeryLongUnicode(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "12€456€7890abcd€f\n€234567890abcde\n€12345"
	expect := "12€456€78…\n€23456789…\n€12345"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}

}

func TestMultipleLong(t *testing.T) {
	var sb strings.Builder
	tw := NewTruncateWriter(&sb, 10)
	str := "1234567890abcdef"
	expect := "123456789…"
	l := len(str)
	n, err := tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result := sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}
	//
	n, err = tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result = sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}
	//
	str = "12345\n\rabcdef"
	expect = "123456789…\n\rabcdef"
	l = len(str)
	n, err = tw.Write([]byte(str))
	if n != l || err != nil {
		t.Errorf(`Write([]byte(%q)) = %v, %v, want %v, nil`, str, n, err, l)
	}
	result = sb.String()
	if result != expect {
		t.Errorf(`sb.String() = %q,  want %q`, result, expect)
	}
}
