package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/xdg"
	"gopkg.in/yaml.v2"
)

// Config contains the root configuration of a project
type Config struct {
	RepoRoot string `yaml:"repo_root"`
}

// LoadConfig checks for configuration files and loads them
func LoadConfig() (c Config) {
	root := xdg.Paths{XDGSuffix: "sagacity"}
	file, err := root.ConfigFile("sagacity.yml")
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Reading file failed: ", file)
	}

	c = Config{}
	yaml.Unmarshal(data, &c)

	return
}
