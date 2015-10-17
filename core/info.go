package core

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	ID   string
	Type string `yaml:"type"`
}

// LoadInfo loads an Info object from a file path
func LoadInfo(p string) (i Info, err error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Reading file failed: ", p)
	}

	basename := filepath.Base(p)
	id := strings.TrimSuffix(basename, filepath.Ext(basename))

	i = Info{ID: id}
	yaml.Unmarshal(data, &i)
	return
}
