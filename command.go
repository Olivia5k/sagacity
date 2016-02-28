package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"os"
	"sort"
)

func commandHostKey(hosts map[string]string) []string {
	ret := make([]string, 0, len(hosts))
	for _, host := range hosts {
		ret = append(ret, host)
	}

	sort.Strings(ret)
	return ret
}

// Command is a representation of an executable command
type Command struct {
	RawType    string            `yaml:"type"`
	RawSummary string            `yaml:"summary"`
	RawCommand string            `yaml:"command"`
	Hosts      map[string]string `yaml:"hosts"`
	id         string
	path       string
	repo       *Repo
}

// MakeCLI creates the CLI tree for a Command info
func (c Command) MakeCLI() []cli.Command {
	sc := make([]cli.Command, 0, len(c.Hosts))
	for _, key := range commandHostKey(c.Hosts) {
		cc := cli.Command{
			Name:     key,
			HideHelp: true,
			Action:   c.Execute,
		}
		sc = append(sc, cc)
	}
	return sc
}

// Execute will execute the command specified by the item.
//
// If the `host` attribute is set, the command will be executed on the host(s)
// specified.
func (c *Command) Execute(cl *cli.Context) {
	blue := color.New(color.FgBlue, color.Bold).SprintfFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintfFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintfFunc()
	green := color.New(color.FgGreen, color.Bold).SprintfFunc()

	args := cl.Args()
	if len(args) == 0 {
		fmt.Println("Specify host targets:")
		for key, def := range c.Hosts {
			fmt.Println(
				fmt.Sprintf(
					"  %s: %s",
					green(key),
					yellow(def),
				),
			)
		}
		return
	}

	hostdef := c.Hosts[args[0]]

	fmt.Println(
		fmt.Sprintf("%s: %s\nRuns %s on hosts matching %s\n",
			blue(c.ID()),
			magenta(c.Summary()),
			yellow(c.RawCommand),
			green(hostdef),
		),
	)

	if !ask("Do you want to continue? [y/N] ") {
		fmt.Println("Doing nothing.")

		os.Exit(1)
	}

	repo := c.repo.ParentRepo()
	host := repo.GetHost(hostdef)
	host.Execute(c.RawCommand)
	return
}

// ID returns the ID of the item
func (c Command) ID() string {
	return c.id
}

// Type returns the Type of the item
func (c Command) Type() string {
	return c.RawType
}

func (c Command) String() string {
	return fmt.Sprintf("C: %s", c.ID())
}

// Path returns the path of the item
func (c Command) Path() string {
	return c.path
}

// Summary returns the summary of the item
func (c Command) Summary() string {
	// TODO(thiderman): This doesn't feel right...
	return c.RawSummary
}

func (c *Command) getHosts(args cli.Args) (names []string) {
	return
}
