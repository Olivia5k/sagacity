package main

import (
	"os"
)

func main() {
	conf := LoadConfig()
	repos := LoadRepositories(conf.RepoRoot)
	app := BuildCLI(repos)
	app.Run(os.Args)
}
