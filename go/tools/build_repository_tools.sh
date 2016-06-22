#!/bin/sh -eu

# Copyright 2016 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Helper script to build //go/tools/gazel in the analysis phase.

GAZEL_PREFIX=github.com/bazelbuild/rules_go/go/tools/gazel
BUILDIFIER_PREFIX=github.com/bazelbuild/buildifier

function prepare_workdir() {
  local dir=`mktemp -d /tmp/gazelXXXXXX`
  if [ $? -ne 0 -o -z "$dir" ]; then
    echo "Failed to create a temporary directory" >&2
    exit 1
  fi

  mkdir -p $dir/pkg/$GAZEL_PREFIX \
    $dir/pkg/$BUILDIFIER_PREFIX \
    $dir/src/$BUILDIFIER_PREFIX/core
  echo $dir
}

function srcs() {
  ls $1/*.go | grep -v '_test\.go$'
}

function go_compile() {
  local go_tool=$1
  local srcdir=$2
  local prefix=$3
  local pkg=$4
  local extra_srcs=${5:-""}
  $go_tool tool compile -pack -I $GOPATH/pkg \
    -o $GOPATH/pkg/$prefix/$pkg.a `srcs $srcdir/$pkg` $extra_srcs
}

function go_yacc() {
  local go_tool=$1
  $go_tool tool yacc -o $3 $2
}

function go_link() {
  local go_tool=$1
  $1 tool link -L $GOPATH/pkg -o $3 $GOPATH/pkg/$2.a
}

function build_buildifier() {
  local go_tool=$1
  local srcdir=$2
  go_yacc $go_tool $srcdir/core/parse.y \
    $GOPATH/src/$BUILDIFIER_PREFIX/core/parse.y.go
  go_compile $go_tool $srcdir $BUILDIFIER_PREFIX core \
    $GOPATH/src/$BUILDIFIER_PREFIX/core/parse.y.go
  go_compile $go_tool $srcdir $BUILDIFIER_PREFIX differ
}

function build_gazel() {
  local go_tool=$1
  local srcdir=$2
  local out=$3
  go_compile $go_tool $srcdir $GAZEL_PREFIX packages
  go_compile $go_tool $srcdir $GAZEL_PREFIX rules
  go_compile $go_tool $srcdir $GAZEL_PREFIX generator
  go_compile $go_tool $srcdir $GAZEL_PREFIX gazel
  go_link $go_tool $GAZEL_PREFIX/gazel $out
}

function main() {
  if [ "$#" -lt 4 ]; then 
    echo "Usage: build_repository_tools.sh /path/to/go_tool " \
      "/path/to/gazel/src /path/to/buldifier/src /path/to/output" >&2
    exit 1
  fi

  GOPATH=`prepare_workdir`
  trap "rm -rf $GOPATH" EXIT
  trap "rm -rf $GOPATH; exit 1" INT TERM

  build_buildifier $1 $3
  build_gazel $1 $2 $4
}

main $@
