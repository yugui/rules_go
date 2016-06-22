/* Copyright 2016 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rules

import (
	"reflect"
	"testing"
)

func TestStructuredResolver(t *testing.T) {
	r := structuredResolver{goPrefix: "example.com/repo"}
	for _, spec := range []struct {
		importpath string
		curPkg     string
		want       label
	}{
		{
			importpath: "example.com/repo",
			curPkg:     "",
			want:       label{name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/lib",
			curPkg:     "",
			want:       label{pkg: "lib", name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/another",
			curPkg:     "",
			want:       label{pkg: "another", name: "go_default_library"},
		},

		{
			importpath: "example.com/repo",
			curPkg:     "lib",
			want:       label{name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/lib",
			curPkg:     "lib",
			want:       label{name: "go_default_library", relative: true},
		},
		{
			importpath: "example.com/repo/lib/sub",
			curPkg:     "lib",
			want:       label{pkg: "lib/sub", name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/another",
			curPkg:     "lib",
			want:       label{pkg: "another", name: "go_default_library"},
		},
	} {

		l, err := r.resolve(spec.importpath, spec.curPkg)
		if err != nil {
			t.Errorf("r.resolve(%q) failed with %v; want success", spec.importpath, err)
			continue
		}
		if got, want := l, spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("r.resolve(%q) = %s; want %s", spec.importpath, got, want)
		}
	}
}

func TestStructuredResolverError(t *testing.T) {
	r := structuredResolver{goPrefix: "example.com/repo"}

	for _, importpath := range []string{
		"example.com/another",
		"example.com/another/sub",
		"example.com/repo_suffix",
	} {
		l, err := r.resolve(importpath, "")
		if err == nil {
			t.Errorf("r.resolve(%q) = %s; want error", importpath, l)
		}
	}
}
