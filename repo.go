package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Repo represents a repository of information yaml files.
type Repo struct {
	Key     string
	Info    map[string]Info
	Control map[string]Info
	root    string
	wg      sync.WaitGroup
}

func (r Repo) String() string {
	return fmt.Sprintf("R: %s (%d articles)", r.Key, len(r.Info))
}

// LoadRepos loads multiple repositories and stores them
func LoadRepos(p string) (repos map[string]Repo) {
	repos = make(map[string]Repo)
	p = getPath(p)
	wg := sync.WaitGroup{}

	files, _ := ioutil.ReadDir(p)
	for _, file := range files {
		fn := filepath.Join(p, file.Name())

		if _, err := os.Stat(filepath.Join(fn, "_repo.yml")); os.IsNotExist(err) {
			log.Println(fmt.Sprintf("Skipping repo %s: no _repo.yml found.", file.Name()))
			continue
		}

		wg.Add(1)
		go func(repos map[string]Repo, fn string, wg *sync.WaitGroup) {
			defer wg.Done()
			repos[asKey(fn)] = NewRepo(fn)
		}(repos, fn, &wg)
	}

	wg.Wait()
	return
}

// NewRepo loads a repository on a path
func NewRepo(p string) (r Repo) {
	r = Repo{Key: asKey(p), root: getPath(p)}
	r.Info = make(map[string]Info)
	r.Control = make(map[string]Info)

	filepath.Walk(r.root, r.walk)
	r.wg.Wait()

	return
}

// Keys returns a sorted list of the info keys in the repository
func (r *Repo) Keys() []string {
	keys := make([]string, 0, len(r.Info))
	for _, info := range r.Info {
		keys = append(keys, info.ID)
	}

	sort.Strings(keys)

	return keys
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
