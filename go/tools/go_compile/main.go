// Command go_compile is a compiler driver similar to a subset of "go build".
// The major differences from "go build" are: first, it provides only
// compilation and cgo support. It does not support linking; second, it takes
// a list of source files rather than a package name.
package main

// NOTE: To keep the bootstraping process simple, this command must consists of
// a single main package.
// It must not depend on anything other than the standard packages.

import (
	"flag"
	"log"
	"path/filepath"
	"strings"
)

var (
	cgo        = flag.Bool("cgo", false, "enable cgo")
	goTool     = flag.String("go-tool", "go", `path to go tool command`)
	outlib     = flag.String("o", "", "output path")
	srcDir     = flag.String("src-dir", ".", "path to the source tree")
	importMaps importMapList
)

func init() {
	flag.Var(&importMaps, "importmap", "same as -importmap in go tool compile")
}

type importMapList []string

func (l *importMapList) Get() interface{} {
	return l
}

func (l *importMapList) String() string {
	return strings.Join([]string(*l), ",")
}

func (l *importMapList) Set(v string) error {
	*l = append(*l, v)
	return nil
}

func run(srcs []string) error {
	if *outlib == "" {
		log.Fatalf("must specify -o flag")
	}
	var err error
	if *outlib, err = filepath.Abs(*outlib); err != nil {
		return err
	}
	if *srcDir, err = filepath.Abs(*srcDir); err != nil {
		return err
	}

	return compile(srcs)
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		log.Fatalln(err)
	}
}
