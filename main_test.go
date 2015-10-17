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
