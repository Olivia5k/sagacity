package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

// BuildCLI builds the base CLI App() object
func BuildCLI(repos map[string]Repo, conf Config) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "sagacity"
	app.Usage = "spread and use knowledge!"

	app.Commands = []cli.Command{
		{
			Name:    "repo",
			Aliases: []string{"r"},
			Usage:   "repo commands",
			Subcommands: []cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "add new repositories",
					Action: func(c *cli.Context) {
						args := c.Args()
						AddRepo(conf.RepoRoot, args[0], args[1])
					},
				},
				{
					Name:    "update",
					Aliases: []string{"u"},
					Usage:   "update repositories",
					Action: func(c *cli.Context) {
						UpdateRepos(repos)
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) {
		// No arguments - print a sorted list of repositories
		args := c.Args()
		if len(args) == 0 {
			ListRepos(repos)
			return
		}

		// One argument - list the items inside the repository
		if len(args) == 1 {
			repo := repos[args[0]]
			repo.Execute()
		}

		// Two arguments - execute the item
		if len(args) == 2 {
			repo := repos[args[0]]
			info := repo.Info[args[1]]
			info.Execute()
		}
	}
	return
}
