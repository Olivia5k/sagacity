package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"os"
)

// Command is a representation of an executable command
type Command struct {
	ID      string
	Summary string `yaml:"summary"`
	Command string `yaml:"command"`
	// TODO(thiderman): Extend to contain Summary
	Hosts map[string]string `yaml:"hosts"`
}

// Execute will execute the command specified by the item.
//
// If the `host` attribute is set, the command will be executed on the host(s)
// specified.
func (c *Command) Execute(r *Repo, cl cli.Args) {
	blue := color.New(color.FgBlue, color.Bold).SprintfFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintfFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintfFunc()
	green := color.New(color.FgGreen, color.Bold).SprintfFunc()

	if len(cl) == 0 {
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

	hostdef := c.Hosts[cl[0]]

	fmt.Println(
		fmt.Sprintf("%s: %s\nRuns %s on hosts matching %s\n",
			blue(c.ID),
			magenta(c.Summary),
			yellow(c.Command),
			green(hostdef),
		),
	)

	if !ask("Do you want to continue? [y/N] ") {
		fmt.Println("Doing nothing.")

		os.Exit(1)
	}

	repo := r.ParentRepo()
	host := repo.GetHost(hostdef)
	host.Execute(c.Command)
	return

	// sh, _ := exec.LookPath("sh")
	// args := []string{sh, "-c", c.Command}

	// cmd := exec.Cmd{
	// 	Path:   sh,
	// 	Args:   args,
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stderr,
	// }
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal("oh noes :(")
	// }

}

func (c *Command) getHosts(args cli.Args) (names []string) {
	return
}
