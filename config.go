package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"os/user"
	"path/filepath"
)

// Config contains the root configuration of a project
type Config struct {
	RepoRoot     string   `yaml:"repository_root"`
	Repositories []string `yaml:"repositories"`
	filename     string
}

// LoadConfig checks for configuration files and loads them
//
// If there is no configuration file, some sane defaults will be provided.
func LoadConfig(fn string) *Config {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		// No configuration file was found - populate the config with defaults

		// Grab the user so we can find the home directory
		u, _ := user.Current()
		root := filepath.Join(u.HomeDir, ".local", "share", "sagacity")
		return &Config{
			RepoRoot:     root,
			Repositories: []string{},
			filename:     fn,
		}
	}

	c := Config{filename: fn}
	yaml.Unmarshal(data, &c)

	return &c
}
