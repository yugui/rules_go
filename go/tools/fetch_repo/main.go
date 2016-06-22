// Command fetch_repo is similar to "go get -d" but it works
// even if the given remote URL is a buildable Go package and
// it checks out a specific revision rather than the latest revision.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/vcs"
)

var (
	remote = flag.String("remote", "", "Go importpath to the repository fetch")
	rev    = flag.String("rev", "", "target reivision")
	dest   = flag.String("dest", "", "destination directory")
)

func create(dir string) error {
	repo, err := vcs.RepoRootForImportPath(*remote, true)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return err
	}
	return repo.VCS.CreateAtRev(dir, repo.Repo, *rev)
}

func update(dir string) error {
	cmd, _, err := vcs.FromDir(dir, filepath.Dir(dir))
	if err != nil {
		return err
	}
	return cmd.TagSync(dir, *rev)
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
