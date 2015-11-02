package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
