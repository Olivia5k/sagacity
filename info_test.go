package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	repos map[string]*Repo
	repo  *Repo
	ctx   *cli.Context
)

// func TestMain(m *testing.M) {
// 	repos = LoadRepos("test/")
// 	d := repos["data"]
// 	ctx = &cli.Context{}
// 	repo = d

// 	retCode := m.Run()

// 	os.Exit(retCode)
// }

func loadTestFile(fn string) Item {
	fn = fmt.Sprintf("test/data/%s.yml", fn)
	i, err := LoadItem(repo, fn)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func TestLoadInfo(t *testing.T) {
	assert := assert.New(t)
	fn := "test/data/first.yaml"
	i, err := LoadItem(&Repo{}, fn)

	assert.Nil(err)
	assert.Equal(i.ID(), "first")
	assert.Equal(i.Type(), "info")
}

// // Executing an info item is just supposed to print the contents.
// func ExampleExecuteInfo() {
// 	i := loadTestFile("ExecuteInfo")
// 	i.Execute(repo, ctx)
// 	// Output: ExecuteInfo content
// }

// TODO(thiderman): Fix this by implementing it
// // Executing a command item is harder to test. The yaml file is just set to echo something.
// func ExampleExecuteCommand() {
// 	i := loadTestFile("ExecuteCommand")
// 	i.Execute(repo, ctx)
// 	// Output: Should there be a 4chan ipsum?
// }
