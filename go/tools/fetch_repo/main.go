// Command fetch_repo is similar to "go get -d" but it work even if the given
// repository path is not a buildable Go package and it checks out a specific
// revision rather than the latest revision.
//
// The difference between fetch_repo and "git clone" or {new_,}git_repository is
// that fetch_repo recognizes import redirection of Go and it supports other
// version control systems.
//
// These differences help us to manage external Go repositories in the manner of
// Bazel.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	remote = flag.String("remote", "", "Go importpath to the repository fetch")
	vcsCmd = flag.String("vcs", "", "vcs")
	rev    = flag.String("rev", "", "target revision")
	dest   = flag.String("dest", "", "destination directory")
)

func create(dir string) error {
	repo, err := fromCmd(*vcsCmd, *remote, dir)
	if err != nil {
		return err
	}
	if err := repo.Get(); err != nil {
		return err
	}
	return repo.UpdateVersion(*rev)
}

func update(dir string) error {
	repo, err := fromDir(dir)
	if err != nil {
		return err
	}
	return repo.UpdateVersion(*rev)
}

func run() error {
	dir, err := filepath.Abs(*dest)
	if err != nil {
		return err
	}
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return create(dir)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}
	return update(dir)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
