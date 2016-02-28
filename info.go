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
	Type() string
	Path() string
	Summary() string
}

// LoadItem loads an Info object from a file path
func LoadItem(r *Repo, p string) (Item, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Reading file failed: ", p)
	}

	// TODO(thiderman): Avoid the double unmarshal.
	// Is there a way we can know some of the data in the stream before the unmarshal?
	// Probably not?
	i := &Info{id: asKey(p), path: p, repo: r}
	yaml.Unmarshal(data, &i)

	switch i.Type() {
	case "command":
		c := &Command{id: asKey(p), path: p, repo: r}
		yaml.Unmarshal(data, &c)
		return c, nil

	case "host":
		h := &HostInfo{id: asKey(p), path: p, repo: r}
		yaml.Unmarshal(data, &h)
		return h, nil
	}

	yaml.Unmarshal(data, &i)
	return i, nil
}

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	RawType    string `yaml:"type"`
	RawSummary string `yaml:"summary"`
	Body       string `yaml:"body"`
	id         string
	path       string
	repo       *Repo
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID())
}

// Execute will print the body
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

// Type returns the type of the item
func (i Info) Type() string {
	return i.RawType
}

// Path returns the path of the item
func (i Info) Path() string {
	return i.path
}

// Summary returns the summary of the item
func (i Info) Summary() string {
	return i.RawSummary
}
