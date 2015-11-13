package main

import (
	"fmt"
	"github.com/fatih/color"
	text "github.com/tonnerre/golang-text"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// ListTypes returns a list of the types in the category map
func ListTypes(c map[string]Category) (keys []string) {
	for key := range c {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// GetHosts returns an array of all the hosts in the category map
func GetHosts(c map[string]Category) (hosts []Host) {
	for _, cat := range c {
		for _, host := range cat.Hosts {
			hosts = append(hosts, host)
		}
	}

	return hosts
}

// PrimaryHost returns the primary host of the category
func PrimaryHost(c map[string]Category) (h *Host) {
	for _, cat := range c {
		if cat.Primary {
			return cat.PrimaryHost()
		}
	}
	return
}

// PrintHost prints a pretty list of the types of hosts
func PrintHost(c map[string]Category) {
	blue := color.New(color.FgBlue, color.Bold).SprintfFunc()
	green := color.New(color.FgGreen, color.Bold).SprintfFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintfFunc()
	yellow := color.New(color.FgYellow).SprintfFunc()
	hiyellow := color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	grey := color.New(color.FgWhite).SprintfFunc()

	for _, t := range ListTypes(c) {
		fmt.Println(fmt.Sprintf("%s:", cyan(t)))
		cat := c[t]
		fmt.Printf("  %s\n", text.Wrap(cat.Summary, 80))
		for x, host := range cat.Hosts {
			// Print the main host item
			fmt.Printf(
				"  %s%s%s %s",
				yellow("["),
				hiyellow(strconv.Itoa(x)),
				yellow("]"),
				blue(host.FQDN),
			)

			// If the host is primary, mark that clearly
			if host.Primary {
				fmt.Printf(" (%s)", green("primary"))
			}

			// If the host has a summary, add that as well
			if host.Summary != "" {
				fmt.Printf(" (%s)", grey(host.Summary))
			}

			fmt.Println()
		}
		fmt.Println()
	}
}

// Category defines a set category of machines
type Category struct {
	Summary string `yaml:"summary"`
	Primary bool   `yaml:"primary"`
	Hosts   []Host `yaml:"hosts"`
}

// PrimaryHost returns the primary host inside of the HostInfo
func (c *Category) PrimaryHost() (h *Host) {
	for _, host := range c.Hosts {
		if host.Primary {
			h = &host
			return
		}
	}

	// No primary was found, just pick the first one
	return &c.Hosts[0]
}

// GetHost returns a specific host, based on FQDN
func (c *Category) GetHost(fqdn string) (h *Host) {
	for _, host := range c.Hosts {
		if host.FQDN == fqdn {
			return &host
		}
	}
	return
}

// Host is a representation of one host
type Host struct {
	FQDN    string `yaml:"fqdn"`
	Summary string `yaml:"summary"`
	Type    string `yaml:"type"`
	Primary bool   `yaml:"primary"`
}

// hasHost returns true if there is a Host definition and false if not.
func (h *Host) hasHost() bool {
	return h.FQDN != ""
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
		log.Fatal("ssh command failed: ", err)
	}
}
