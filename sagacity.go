package main

import (
	"fmt"
	"os"
)

func main() {
	conf := LoadConfig()
	fmt.Println(conf)

	NewRepo(conf.RepoRoot)

	app := BuildCLI()
	app.Run(os.Args)
}
