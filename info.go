package main

import (
	"fmt"
	text "github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	ID      string
	Type    string `yaml:"type"`
	Summary string `yaml:"summary"`
	Body    string `yaml:"body"`
	Host    string `yaml:"host"`
	Command string `yaml:"command"`
}

func (i Info) String() string {
	return fmt.Sprintf("I: %s", i.ID)
}

// Execute will figure out the type of the info and execute accordingly
func (i *Info) Execute() {
	if i.Type == "info" {
		i.PrintBody()
	} else if i.Type == "command" {
		i.ExecuteCommand()
	}
}

// PrintBody will pretty format the body of the item
func (i *Info) PrintBody() {
	out := text.Wrap(i.Body, 80)
	fmt.Println(out)
}

// ExecuteCommand will execute the command specified by the item
func (i *Info) ExecuteCommand() {
	sh, _ := exec.LookPath("sh")
	args := []string{sh, "-c", i.Command}

	cmd := exec.Cmd{
		Path:   sh,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal("oh noes :(")
	}
}
