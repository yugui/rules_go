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

package generator

import (
	"fmt"
	"go/build"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/bazelbuild/rules_go/go/tools/gazel/testdata"
)

func TestGenerator(t *testing.T) {
	stub := stubRuleGen{
		fixtures: map[string][]*bzl.Rule{
			"lib": {
				{
					Call: &bzl.CallExpr{
						X: &bzl.LiteralExpr{Token: "go_library"},
					},
				},
			},
			"bin": {
				{
					Call: &bzl.CallExpr{
						X: &bzl.LiteralExpr{Token: "go_binary"},
					},
				},
			},
		},
	}

	repo := filepath.Join(testdata.Dir(), "repo")
	g, err := New(repo)
	if err != nil {
		t.Errorf("New(%q) failed with %v; want success", repo, err)
		return
	}
	g.g = stub

	got, err := g.Generate(repo)
	if err != nil {
		t.Errorf("g.Generate(%q) failed with %v; want success", repo, err)
	}
	sort.Sort(fileSlice(got))

	want := []*bzl.File{
		{
			Path: "lib/BUILD",
			Stmt: []bzl.Expr{stub.fixtures["lib"][0].Call},
		},
		{
			Path: "bin/BUILD",
			Stmt: []bzl.Expr{stub.fixtures["bin"][0].Call},
		},
	}
	sort.Sort(fileSlice(want))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("g.Generate(%q) = %v; want %v", repo, prettyFiles(got), prettyFiles(want))
	}
}

type prettyFiles []*bzl.File

func (p prettyFiles) String() string {
	var items []string
	for _, f := range p {
		items = append(items, fmt.Sprintf("{Path: %q, Stmt: %q", f.Path, string(bzl.Format(f))))
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ","))
}

type fileSlice []*bzl.File

func (p fileSlice) Less(i, j int) bool { return strings.Compare(p[i].Path, p[j].Path) < 0 }
func (p fileSlice) Len() int           { return len(p) }
func (p fileSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// stubRuleGen is a test stub implementation of rules.Generator
type stubRuleGen struct {
	fixtures map[string][]*bzl.Rule
}

func (s stubRuleGen) Generate(rel string, pkg *build.Package) ([]*bzl.Rule, error) {
	return s.fixtures[rel], nil
}
