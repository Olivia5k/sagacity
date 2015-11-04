package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// Host is a representation one or more hosts
type Host struct {
	FQDN string `yaml:"fqdn"`
	Role string `yaml:"role"`
}

func (h *Host) getHost() string {
	if h.FQDN == "" {
		return "localhost"
	}
	return h.FQDN
}

// hasHost returns true if there is a Host definition and false if not.
func (h *Host) hasHost() bool {
	return h.FQDN != "" || h.Role != ""
}

// Execute runs a command on the server
// The default is to open a shell. If arguments are given, those arguments
// will be executed verbatim on the host.
func (h *Host) Execute(extra string) {
	ssh, _ := exec.LookPath("ssh")

	// Split the extra arguments into an array...
	e := strings.Split(extra, " ")

	// ...and join them into the arguments list.
	args := append([]string{ssh, h.FQDN, "-A", "-t"}, e...)

	cmd := exec.Cmd{
		Path:   ssh,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal("oh noes :(")
	}
}
