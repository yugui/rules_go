package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func compile(srcs []string) error {
	cmd := exec.Command(*goTool, "tool", "compile")
	cmd.Args = append(
		cmd.Args,
		"-o", *outlib,
		"-pack",
		"-I", *srcDir,
		"-trimpath", *srcDir,
	)
	for _, m := range importMaps {
		cmd.Args = append(cmd.Args, "-importmap", m)
	}
	for _, s := range srcs {
		if !filepath.IsAbs(s) {
			s = filepath.Join(*srcDir, s)
		}
		cmd.Args = append(cmd.Args, s)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
