package main

import (
	"fmt"
	"os"

	"github.com/thiderman/sagacity/core"
)

func main() {
	conf := core.LoadConfig()
	fmt.Println(conf)

	core.NewRepo(conf.RepoRoot)

	app := core.BuildCLI()
	app.Run(os.Args)
}
