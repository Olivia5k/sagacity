package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

func ask(prompt string) bool {
	var resp string

	fmt.Print(prompt)
	_, err := fmt.Scanln(&resp)
	if err != nil {
		if err.Error() != "unexpected newline" && err.Error() != "EOF" {
			log.Fatal(err)
		} else {
			return false
		}
	}

	if string(strings.ToLower(resp)[0]) == "y" {
		return true
	}
	return false
}
