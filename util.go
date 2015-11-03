package main

import (
	"log"
	"os"
	"os/exec"
)

// Helper for executing git commands
func git(pwd string, args ...string) {
	if pwd == "" {
		pwd, _ = os.Getwd()
	}
	git, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("no git :'(   ", err)
	}

	// why................
	args = append([]string{git}, args...)

	cmd := exec.Cmd{
		Path:   git,
		Args:   args,
		Dir:    pwd,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
	if err != nil {
		log.Println("git command failed - aborting")
		log.Fatal(err)
	}
}
