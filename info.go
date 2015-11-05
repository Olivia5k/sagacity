package main

import (
	"fmt"
	"github.com/fatih/color"
	text "github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	Type    string `yaml:"type"`
	Summary string `yaml:"summary"`
	Body    string `yaml:"body"`
	Command string `yaml:"command"`
	Hosts   []Host `yaml:"hosts"`
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID)
}

// Execute will figure out the type of the info and execute accordingly
func (i *Info) Execute() {
	if i.Type == "info" {
		i.PrintBody()
	} else if i.Type == "command" {
		i.ExecuteCommand()
	} else if i.Type == "host" {
		i.ExecuteHost()
	} else {
		log.Fatal("Unknown type:", i.Type)
	}
}

// PrintBody will pretty format the body of the item
func (i *Info) PrintBody() {
	out := text.Wrap(i.Body, 80)
	fmt.Println(out)
}

// ExecuteCommand will execute the command specified by the item.
//
// If the `host` attribute is set, the command will be executed on the host(s)
// specified.
func (i *Info) ExecuteCommand() {
	blue := color.New(color.FgBlue, color.Bold).SprintfFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintfFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintfFunc()
	green := color.New(color.FgGreen, color.Bold).SprintfFunc()

	fmt.Println(
		fmt.Sprintf("%s: %s\nRuns %s on %s\n",
			blue(i.ID),
			magenta(i.Summary),
			yellow(i.Command),
			green(strings.Join(i.getHosts(), ", ")),
		),
	)

	if !ask("Do you want to continue? [y/N] ") {
		fmt.Println("Doing nothing.")
		os.Exit(1)
	}

	host := i.GetHost()
	if host.hasHost() {
		host.Execute(i.Command)
		return
	}

	sh, _ := exec.LookPath("sh")
	args := []string{sh, "-c", i.Command}

	cmd := exec.Cmd{
		Path:   sh,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal("oh noes :(")
	}

}

// ExecuteHost opens a ssh connection to the specified host
func (i *Info) ExecuteHost() {
	i.GetHost().Execute("") // Called with no args - new ssh session
}

// GetHost will return the primary host of the item
func (i *Info) GetHost() *Host {
	return &i.Hosts[0]
}

func (i *Info) getHosts() []string {
	hosts := make([]string, 0, len(i.Hosts))
	for _, host := range i.Hosts {
		if host.FQDN == "" {
			hosts = append(hosts, "localhost")
		} else {
			hosts = append(hosts, host.FQDN)
		}
	}

	return hosts
}
