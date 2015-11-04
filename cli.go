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

		// More arguments.
		// If the first argument given matches a repo, use that one and run
		// repo.Execute() with the cli context so it can determine what to do.
		//
		// If not, list the available repositories.
		if repo, ok := repos[args[0]]; ok {
			repo.Execute(c)
		} else {
			fmt.Printf("%s: no such repo. Available repos are:\n\n", args[0])
			ListRepos(repos)
		}
	}
	return
}
