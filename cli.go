package main

import (
	"github.com/codegangsta/cli"
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

	// Repo management commands are always present.
	commands = append(commands, []cli.Command{
		{
			Name:     "repo",
			Usage:    "repo commands",
			HideHelp: true,
			Subcommands: []cli.Command{
				{
					Name:     "add",
					Usage:    "add new repositories",
					HideHelp: true,
					Action: func(c *cli.Context) {
						args := c.Args()
						AddRepo(conf.RepoRoot, args[0], args[1])
					},
				},
				{
					Name:     "update",
					Usage:    "update repositories",
					HideHelp: true,
					Action: func(c *cli.Context) {
						UpdateRepos(repos)
					},
				},
			},
		},
	}...)

	app.Commands = commands
	return
}
