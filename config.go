package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/xdg"
	"gopkg.in/yaml.v2"
)

func getConfigFilename() string {
	root := xdg.Paths{XDGSuffix: "sagacity"}
	file, err := root.ConfigFile("sagacity.yaml")
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// Config contains the root configuration of a project
type Config struct {
	RepoRoot string `yaml:"repo_root"`
}

// LoadConfig checks for configuration files and loads them
func LoadConfig() (c Config) {
	data, err := ioutil.ReadFile(getConfigFilename())
	if err != nil {
		log.Fatal("No root configuration file found: ", err)
	}

	c = Config{}
	yaml.Unmarshal(data, &c)

	return
}
