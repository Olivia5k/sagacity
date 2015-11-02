package main

import (
	"fmt"
	"os"
)

func main() {
	conf := LoadConfig()
	repos := LoadRepositories(conf.RepoRoot)

	fmt.Println(repos)

	app := BuildCLI()
	app.Run(os.Args)
}
