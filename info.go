package main

import (
	"fmt"
	text "github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getPath(p string) string {
	path, _ := filepath.Abs(p)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(err)
	}
	return path
}

func asKey(p string) string {
	basename := filepath.Base(p)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	ID      string
	Type    string `yaml:"type"`
	Summary string `yaml:"summary"`
	Body    string `yaml:"body"`
}

// LoadInfo loads an Info object from a file path
func LoadInfo(p string) (i Info, err error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Reading file failed: ", p)
	}

	i = Info{ID: asKey(p)}
	yaml.Unmarshal(data, &i)
	return
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID)
}

func (i *Info) Execute() {
	if i.Type == "info" {
		i.PrintBody()
	}
}

// PrintBody will pretty format the body of the item
func (i *Info) PrintBody() {
	out := text.Wrap(i.Body, 80)
	fmt.Println(out)
}
