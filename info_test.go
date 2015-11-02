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

func TestLoadInfo(t *testing.T) {
	assert := assert.New(t)
	fn := "test/data/first.yml"
	i, err := LoadInfo(fn)

	assert.Nil(err)
	assert.Equal(i.ID, "first")
	assert.Equal(i.Type, "info")
}

func TestNewRepoLoadsTheFirstFile(t *testing.T) {
	assert := assert.New(t)
	dir := "test/data/"
	r := NewRepo(dir)

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
