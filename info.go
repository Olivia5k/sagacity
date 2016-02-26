package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// An Item is a representation of the YAML files in the repositories
type Item interface {
	Execute(c *cli.Context)
	String() string
	MakeCLI() []cli.Command
	ID() string
	Path() string
	Summary() string
}

// LoadItem loads an Info object from a file path
func LoadItem(r *Repo, p string) (i Item, err error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Reading file failed: ", p)
	}

	i = &Info{id: asKey(p), path: p}
	yaml.Unmarshal(data, &i)
	return
}

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	Type        string `yaml:"type"`
	SummaryText string `yaml:"summary"`
	Body        string `yaml:"body"`
	id          string
	path        string
	repo        *Repo
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID())
}

// Execute will figure out the type of the info and execute accordingly
func (i Info) Execute(c *cli.Context) {
	out := text.Wrap(i.Body, 80)
	fmt.Println(out)
}

// MakeCLI makes a dummy CLI - Info items have no subcommands
func (i Info) MakeCLI() []cli.Command {
	return []cli.Command{}
}

// ID returns the ID of the item
func (i Info) ID() string {
	return i.id
}

// Path returns the path of the item
func (i Info) Path() string {
	return i.path
}

// Summary returns the summary of the item
func (i Info) Summary() string {
	// TODO(thiderman): This doesn't feel right...
	return i.SummaryText
}
