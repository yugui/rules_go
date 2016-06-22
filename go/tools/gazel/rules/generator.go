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
	"go/build"
	"path/filepath"
	"strings"

	bzl "github.com/bazelbuild/buildifier/core"
)

// Generator generates Bazel build rules for Go build targets
type Generator interface {
	// Generate generates build rules for build targets in a Go package in a
	// repository.
	//
	// "rel" is a relative slash-separated path from the repostiry root
	// directory to the Go package directory. It is empty if the package
	// directory is the repository root itself.
	// "pkg" is a description about the package.
	Generate(rel string, pkg *build.Package) ([]*bzl.Rule, error)
}

// NewGenerator returns an implementation of Generator.
//
// "goPrefix" is the go_prefix corresponding to the repository root.
// See also https://github.com/bazelbuild/rules_go#go_prefix.
func NewGenerator(goPrefix string) Generator {
	return &generator{
		// TODO(yugui) Support another resolver to cover the pattern 2 in
		// https://github.com/bazelbuild/rules_go/issues/16#issuecomment-216010843
		r: structuredResolver{goPrefix: goPrefix},
	}
}

type generator struct {
	r labelResolver
}

func (g *generator) Generate(rel string, pkg *build.Package) ([]*bzl.Rule, error) {
	r, err := g.generate(rel, pkg)
	if err != nil {
		return nil, err
	}
	rules := []*bzl.Rule{r}

	if len(pkg.TestGoFiles) > 0 {
		t, err := g.generateTest(rel, pkg, r.AttrString("name"))
		if err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}

	if len(pkg.XTestGoFiles) > 0 {
		t, err := g.generateXTest(rel, pkg, r.AttrString("name"))
		if err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}
	return rules, nil
}

func (g *generator) generate(rel string, pkg *build.Package) (*bzl.Rule, error) {
	kind := "go_library"
	name := "go_default_library"
	if pkg.IsCommand() {
		kind = "go_binary"
		name = filepath.Base(pkg.Dir)
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: pkg.GoFiles},
	}

	deps, err := g.dependencies(pkg.Imports, rel)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}

	return newRule(kind, nil, attrs)
}

func (g *generator) generateTest(rel string, pkg *build.Package, library string) (*bzl.Rule, error) {
	name := library + "_test"
	if library == "go_default_library" {
		name = "go_default_test"
	}
	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: pkg.TestGoFiles},
		{key: "library", value: ":" + library},
	}

	deps, err := g.dependencies(pkg.TestImports, rel)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}
	return newRule("go_test", nil, attrs)
}

func (g *generator) generateXTest(rel string, pkg *build.Package, library string) (*bzl.Rule, error) {
	name := library + "_xtest"
	if library == "go_default_library" {
		name = "go_default_xtest"
	}
	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: pkg.XTestGoFiles},
	}

	deps, err := g.dependencies(pkg.XTestImports, rel)
	if err != nil {
		return nil, err
	}
	attrs = append(attrs, keyvalue{key: "deps", value: deps})
	return newRule("go_test", nil, attrs)
}

func (g *generator) dependencies(imports []string, dir string) ([]string, error) {
	var deps []string
	for _, p := range imports {
		if isStandard(p) {
			continue
		}
		l, err := g.r.resolve(p, dir)
		if err != nil {
			return nil, err
		}
		deps = append(deps, l.String())
	}
	return deps, nil
}

// isStandard determines if importpath points a Go standard package.
func isStandard(importpath string) bool {
	seg := strings.SplitN(importpath, "/", 2)[0]
	return !strings.Contains(seg, ".")
}
