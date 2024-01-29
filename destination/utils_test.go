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

func TestToTitleCase(t *testing.T) {
	is := is.New(t)

	is.Equal(toTitleCase("hello world"), "Hello World")
	is.Equal(toTitleCase("HELLO WORLD"), "Hello World")
	is.Equal(toTitleCase("title Case"), "Title Case")
}

func TestContainsString(t *testing.T) {
	is := is.New(t)

	slice := []string{"Red", "Green", "Blue"}
	is.True(containsString(slice, "Red"))     // should find "Red"
	is.True(!containsString(slice, "Yellow")) // should not find "Yellow"
}
