package main

import (
	"github.com/codegangsta/cli"
	"os"
	"sort"
)

// BuildCLI builds the base CLI App() object
func BuildCLI(repos map[string]*Repo, conf *Config) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "sp"
	app.EnableBashCompletion = true
	app.Usage = "spread and use knowledge!"
	app.HideHelp = true

	repolen := len(repos)
	commands := make([]cli.Command, 0, repolen+2)

	keys := make([]string, 0, repolen)
	for key := range repos {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		repo := repos[key]
		commands = append(commands, repo.MakeCLI())
	}

	// Repo management commands are only present if we are not doing bash completion.
	if !isCompleting() {
		commands = append(commands, []cli.Command{
			{
				Name:     "repo",
				Usage:    "repo commands",
				HideHelp: true,
				Subcommands: []cli.Command{
					{
						Name:     "add",
						Usage:    "add <url>",
						HideHelp: true,
						Action: func(c *cli.Context) {
							args := c.Args()
							AddRepo(conf, args[0])
						},
					},
					{
						Name:     "update",
						Usage:    "update",
						HideHelp: true,
						Action: func(c *cli.Context) {
							UpdateRepos(repos)
						},
					},
				},
			},
		}...)
	}

	app.Commands = commands
	return
}

// isCompleting returns boolean if we are doing bash completion or not
//
// This is only really used by BuildCLI() when determining what to show. To
// not clutter the base command level, we avoid adding the repo management
// commands whenever we are doing completion. Clean!
func isCompleting() bool {
	for _, arg := range os.Args {
		if arg == "--generate-bash-completion" {
			return true
		}

	}
	return false
}
