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
