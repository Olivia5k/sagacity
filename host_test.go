package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestHostInfoHostsAreSet(t *testing.T) {
	assert := assert.New(t)
	p := "test/hosts/db.yaml"
	h := HostInfo{id: asKey(p), path: p, repo: &Repo{}}

	data, _ := ioutil.ReadFile(p)
	yaml.Unmarshal(data, &h)

	t.Log(h.String())

	assert.Equal(5, len(h.Types))
}
