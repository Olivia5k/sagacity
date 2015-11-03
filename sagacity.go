package main

import (
	"os"
)

func main() {
	conf := LoadConfig()
	repos := LoadRepos(conf.RepoRoot)
	app := BuildCLI(repos)
	app.Run(os.Args)
}
