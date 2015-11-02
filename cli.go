package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"sort"
)

// BuildCLI builds the base CLI App() object
func BuildCLI(repos map[string]Repo) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "sagacity"
	app.Usage = "spread and use knowledge!"

	app.Action = func(c *cli.Context) {
		// No arguments - print a sorted list of repositories
		args := c.Args()
		if len(args) == 0 {
			keys := make([]string, 0, len(repos))
			for key := range repos {
				keys = append(keys, key)
			}

			sort.Strings(keys)
			for _, key := range keys {
				fmt.Println(key)
			}

			return
		}

		// One argument - list the items inside the repository
		if len(args) == 1 {
			repo := repos[args[0]]
			for _, key := range repo.Keys() {
				fmt.Println(key)
			}
		}

		// Two arguments - print the item
		if len(args) == 2 {
			repo := repos[args[0]]
			info := repo.Info[args[1]]
			info.PrintBody()
		}
	}
	return
}
