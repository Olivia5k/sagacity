package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"os"
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

// persist saves the file to disk
func (c *Config) persist() error {
	// Create the directory if it doesn't exist
	os.MkdirAll(c.RepoRoot, 0755)

	d, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filename, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AddRepo adds a new repository to the config and saves the YAML
func (c *Config) AddRepo(dir string) error {
	c.Repositories = append(c.Repositories, dir)
	return c.persist()
}
