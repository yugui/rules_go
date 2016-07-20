package main

import (
	"fmt"

	mvcs "github.com/Masterminds/vcs"
	"golang.org/x/tools/go/vcs"
)

type repo struct {
	mvcs.Repo
}

func fromCmd(cmd, importpath, dir string) (mvcs.Repo, error) {
	repo, err := vcs.RepoRootForImportPath(*remote, true)
	if err != nil {
		return nil, err
	}
	if importpath != repo.Root {
		return nil, fmt.Errorf("not a root of a repository: %s", importpath)
	}
	switch repo.VCS.Name {
	case "Git":
		return mvcs.NewGitRepo(repo.Repo, dir)
	case "Mercurial":
		return mvcs.NewHgRepo(repo.Repo, dir)
	case "Subversion":
		return mvcs.NewSvnRepo(repo.Repo, dir)
	case "Bazzar":
		return mvcs.NewBzrRepo(repo.Repo, dir)
	default:
		return nil, fmt.Errorf("unsupported VCS: %s", repo.VCS)
	}
}

func fromDir(dir string) (mvcs.Repo, error) {
	repo, err := vcs.RepoRootForImportPath(*remote, true)
	if err != nil {
		return nil, err
	}
	return mvcs.NewRepo(repo.Repo, dir)
}
