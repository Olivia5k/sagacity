package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// createJunk creates a lot of garbage files in a temporary diretory
func createJunk(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.FileMode(0700))
	}

	files, _ := ioutil.ReadDir(dir)
	for i := len(files); i < 100; i++ {
		// Since we are not using strconv, the strings created will almost
		// literally be junk characters. Rather than numbers, they will be
		// characters and special characters.
		path := filepath.Join(dir, string(i))
		os.Create(path)
	}
}

func TestNewRepoLoadsTheFirstFile(t *testing.T) {
	assert := assert.New(t)
	dir := "test/data/"
	r := NewRepo(dir)

	assert.True(len(r.Info) >= 1)
	assert.Equal(r.Info["first"].ID, "first")
}

func TestNewRepoHandlesControlFiles(t *testing.T) {
	assert := assert.New(t)
	dir := "test/data/"
	r := NewRepo(dir)

	assert.Equal(len(r.Info), 2)
	assert.Equal(len(r.Control), 1)
}

func TestNewRepoLoadsMultipleFiles(t *testing.T) {
	assert := assert.New(t)
	dir := "test/data/"
	r := NewRepo(dir)

	assert.Equal(r.Info["first"].ID, "first")
	assert.Equal(r.Info["second"].ID, "second")
	assert.Equal(len(r.Info), 2)
}

// Generate tons of junk files, and check that none of them are loaded
func TestNewRepoIgnoresNonYamlJunk(t *testing.T) {
	assert := assert.New(t)
	dir := "test/junk/"
	createJunk(dir)
	r := NewRepo(dir)

	assert.Equal(len(r.Info), 0)
	assert.Equal(len(r.Control), 0)
}

func TestNewRepoNestsDeep(t *testing.T) {
	assert := assert.New(t)

	dir := "test/deep/"
	r := NewRepo(dir)

	one := r.Subrepos["one"]
	two := one.Subrepos["two"]
	three := two.Subrepos["three"]
	four := three.Subrepos["four"]
	five := four.Subrepos["five"]

	assert.Equal(five.Info["first"].Body, "glitteringprizes")
}

func TestNewRepoNestsDeepAndDoesNotPutItemsOnTopLevels(t *testing.T) {
	assert := assert.New(t)

	dir := "test/deep/"
	r := NewRepo(dir)

	one := r.Subrepos["one"]
	two := one.Subrepos["two"]
	three := two.Subrepos["three"]
	four := three.Subrepos["four"]
	five := four.Subrepos["five"]

	assert.Equal(len(one.Info), 0)
	assert.Equal(len(two.Info), 0)
	assert.Equal(len(three.Info), 0)
	assert.Equal(len(four.Info), 0)
	assert.Equal(len(five.Info), 1)
}
