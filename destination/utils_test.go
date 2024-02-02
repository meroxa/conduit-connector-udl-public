// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package destination

import (
	"testing"

	"github.com/matryer/is"
)

func TestSupportedStringValues(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		name      string
		check     string
		supported []string
		want      bool
	}{
		{"Test present value", "TEST", DataModeValues, true},
		{"Case insensitive test", "test", DataModeValues, true},
		{"Trim spaces and case", " test ", DataModeValues, true},
		{"Unsupported value", "UNKNOWN", DataModeValues, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t) // New instance for concurrent tests
			result := SupportedStringValues(tc.check, tc.supported)
			is.Equal(result, tc.want) // result should be what we expect
		})
	}
}

func TestAllZero(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		name  string
		slice []float64
		want  bool
	}{
		{"All zeros", []float64{0, 0, 0}, true},
		{"Contains non-zero", []float64{0, 1.2, 0}, false},
		{"Empty slice", []float64{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			result := allZero(tc.slice)
			is.Equal(result, tc.want) // result should be what we expect
		})
	}
}

func TestReplaceUnderscoresWithSpaces(t *testing.T) {
	is := is.New(t)

	is.Equal(replaceUnderscoresWithSpaces("hello_world"), "hello world")
	is.Equal(replaceUnderscoresWithSpaces("_"), " ")
	is.Equal(replaceUnderscoresWithSpaces(""), "")
}

func TestToTitleCase(t *testing.T) {
	is := is.New(t)

	is.Equal(toTitleCase("hello world"), "Hello World")
	is.Equal(toTitleCase("HELLO WORLD"), "Hello World")
	is.Equal(toTitleCase("title Case"), "Title Case")

}
