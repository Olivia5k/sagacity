package main

import (
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	u, _ := user.Current()
	fn := filepath.Join(u.HomeDir, ".config", "sagacity", "sagacity.yaml")
	conf := LoadConfig(fn)

	repos := LoadRepos(conf.RepoRoot)
	app := BuildCLI(repos, conf)
	app.Run(os.Args)
}
