package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

func asKey(p string) string {
	basename := filepath.Base(p)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

// Info is the main storage for information. All yaml files map to this.
type Info struct {
	ID   string
	Type string `yaml:"type"`
	Body string `yaml:"body"`
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

// Repo represents a repository of information yaml files.
type Repo struct {
	Key     string
	Info    map[string]Info
	Control map[string]Info
	root    string
	wg      sync.WaitGroup
}

// NewRepo loads a repository on a path
func NewRepo(p string) (r Repo) {
	path, _ := filepath.Abs(p)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("no such file or directory: %s", path)
	}

	r = Repo{Key: asKey(p), root: path}
	r.Info = make(map[string]Info)
	r.Control = make(map[string]Info)

	filepath.Walk(r.root, r.walk)
	r.wg.Wait()

	return
}

func (r *Repo) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("walk error: ", err)
		return err
	}

	if strings.HasSuffix(path, ".yml") {
		r.wg.Add(1)
		go r.loadInfo(path)
	}

	return nil
}

func (r *Repo) loadInfo(path string) {
	defer r.wg.Done()

	info, err := LoadInfo(path)
	if err != nil {
		log.Println("Failed to load info: ", err)
	}

	// Control files start with an underscore and should not be stored as
	// normal Info documents.
	if strings.HasPrefix(asKey(path), "_") {
		r.Control[info.ID] = info
	} else {
		r.Info[info.ID] = info
	}
}
