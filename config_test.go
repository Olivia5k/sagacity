package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadConfigFileDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	fn := "nonexistant/config.yaml"
	c := LoadConfig(fn)

	assert.Equal(0, len(c.Repositories))
	assert.Equal(fn, c.filename)
	assert.True(strings.HasSuffix(c.RepoRoot, ".local/share/sagacity"))
}

func TestLoadConfigFileDoesExist(t *testing.T) {
	assert := assert.New(t)
	fn := "test/config_load_test.yaml"
	c := LoadConfig(fn)

	assert.Equal(2, len(c.Repositories))
	assert.Equal("/whisky/in/the/jar", c.Repositories[0])
	assert.Equal("/rapunzel/hair", c.Repositories[1])

	assert.Equal(fn, c.filename)
	assert.Equal("/fiddler/on/the/green", c.RepoRoot)
}
