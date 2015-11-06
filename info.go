package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	// "github.com/fatih/color"
	text "github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	// "os/exec"
	"strconv"
	"strings"
)

// LoadInfo loads an Info object from a file path
func LoadInfo(p string) (i Info, err error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Reading file failed: ", p)
	}

	i = Info{ID: asKey(p)}
	yaml.Unmarshal(data, &i)
	return
}

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	ID      string
	Type    string              `yaml:"type"`
	Summary string              `yaml:"summary"`
	Body    string              `yaml:"body"`
	Command Command             `yaml:"command"`
	Hosts   map[string]Category `yaml:"types"`
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID)
}

// Execute will figure out the type of the info and execute accordingly
func (i *Info) Execute(r *Repo, c cli.Args) {
	if i.Type == "info" {
		i.PrintBody()
	} else if i.Type == "command" {
		i.Command.Execute(r, c)
	} else if i.Type == "host" {
		i.ExecuteHost(c)
	} else {
		log.Fatal("Unknown type:", i.Type)
	}
}

// PrintBody will pretty format the body of the item
func (i *Info) PrintBody() {
	out := text.Wrap(i.Body, 80)
	fmt.Println(out)
}

// ExecuteHost opens a ssh connection to the specified host
func (i *Info) ExecuteHost(c cli.Args) {
	clen := len(c)
	switch clen {
	case 0:
		// No further arguments - we have selected a host entry but no type.
		// Print the list of hosts.
		PrintHost(i.Hosts)
	case 1, 2:
		t := c[0]
		if cat, ok := i.Hosts[t]; ok {
			if clen == 1 {
				// One argument, go to the primary of that category
				cat.PrimaryHost().Execute("")
			} else {
				// Two arguments, go to specified host
				x, err := strconv.Atoi(c[1])
				if err != nil {
					log.Fatal("Non-integer argument:", c[1])
				}

				host := cat.Hosts[x]
				host.Execute("")
			}

		} else {
			fmt.Println("No such type:", t)
			fmt.Println(
				fmt.Sprintf("Choices are: %s", strings.Join(ListTypes(i.Hosts), ", ")),
			)
			os.Exit(1)
		}
	}
}

// getHosts gets a string representation of all of the hosts in the item
func (i *Info) getHosts() (hosts []string) {
	for _, host := range GetHosts(i.Hosts) {
		hosts = append(hosts, host.FQDN)
	}

	return
}
