package main

import (
// "testing"
)

func ExampleCliNoArguments() {
	repos := map[string]Repo{
		"zathura": Repo{},
		"test":    Repo{},
		"gamma":   Repo{},
	}

	app := BuildCLI(repos)
	args := make([]string, 1)

	app.Run(args)
	// Output: gamma
	// test
	// zathura
}

func ExampleCliPrintItem() {
	repos := map[string]Repo{
		"joanjett": Repo{
			Info: map[string]Info{
				"bad_reputation": Info{
					Body: "I don't give a damn about my bad reputation!",
				},
			},
		},
	}

	app := BuildCLI(repos)
	args := []string{"/go/bin/sagacity", "joanjett", "bad_reputation"}

	app.Run(args)
	// Output: I don't give a damn about my bad reputation!
}
