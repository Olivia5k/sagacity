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
	"path/filepath"
	"strings"
)

func getPath(p string) string {
	path, _ := filepath.Abs(p)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(err)
	}
	return path
}

func asKey(p string) string {
	basename := filepath.Base(p)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

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
	Host    Host   `yaml:"host"`
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
			green(i.Host.getHost()),
		),
	)

	if !ask("Do you want to continue? [y/N] ") {
		fmt.Println("Doing nothing.")
		os.Exit(1)
	}

	if i.Host.hasHost() {
		i.Host.Execute(i.Command)
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
	i.Host.Execute("") // Called with no args - new ssh session
}
