package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// BuildCLI builds the base CLI App() object
func BuildCLI() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "sagacity"
	app.Usage = "spread and use knowledge!"
	app.Action = func(c *cli.Context) {
		fmt.Println("hehe")
	}

	return
}
