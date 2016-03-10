package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Repo represents a repository of information yaml files.
type Repo struct {
	Key      string `yaml:"key"`
	Summary  string `yaml:"summary"`
	Alias    string `yaml:"alias"`
	Items    map[string]Item
	Control  map[string]Item
	Subrepos map[string]*Repo
	Parent   *Repo
	root     string
}

func (r Repo) String() string {
	return fmt.Sprintf("R: %s (%d articles)", r.Key, len(r.Items))
}

// LoadRepos loads multiple repositories and stores them
func LoadRepos(c *Config) (repos map[string]*Repo) {
	repos = make(map[string]*Repo)
	cr := make(chan *Repo)

	started := 0
	for _, file := range c.Repositories {
		if _, err := os.Stat(filepath.Join(file, "_repo.yaml")); os.IsNotExist(err) {
			// log.Println(fmt.Sprintf("Skipping repo %s: no _repo.yaml found.", file.Name()))
			continue
		}

		started++
		go func(c chan<- *Repo, fn string) {
			c <- NewRepo(fn)
		}(cr, file)
	}

	for x := 0; x < started; x++ {
		r := <-cr
		if r != nil {
			repos[r.Key] = r
		}
	}

	return
}

// UpdateRepos will run git pull on the repos
func UpdateRepos(repos map[string]*Repo) {
	for key, repo := range repos {
		log.Printf("Updating %s...", key)
		repo.git("pull", "origin", "master")
	}
}

// AddRepo clones a new repository
func AddRepo(config *Config, url string) {
	// Clean the name of prefixes and stuff, leaving just the trailing word. This
	// lets us use `saga-topic` or `kb-topic` or whatever and we'll still get
	// just `topic` when we're grabbing.
	rxp := regexp.MustCompile(".*-")
	name := rxp.ReplaceAllString(url, "")

	// Clone the repo! |o/
	dir := filepath.Join(config.RepoRoot, name)
	git("", "clone", url, dir)

	// Persist the changes into the configuration file
	err := config.AddRepo(dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Added %s as %s!\n", url, name)
}

// NewRepo loads a repository on a path
func NewRepo(p string) *Repo {
	var subdirs []string
	var items []string

	p = getPath(p)
	r := Repo{Key: asKey(p), root: p}

	// Check if this is a root repo. If it is, load the data from the _repo.yaml file into
	// the newly created repo.
	rfile := filepath.Join(p, "_repo.yaml")
	if _, err := os.Stat(rfile); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(rfile)

		if err != nil {
			log.Fatal("Reading repo file failed: ", p)
		}
		yaml.Unmarshal(data, &r)
	}

	r.Items = make(map[string]Item)
	r.Control = make(map[string]Item)
	r.Subrepos = make(map[string]*Repo)

	files, _ := ioutil.ReadDir(p)

	// Loop through the files and put files and dirs in different lists
	for _, f := range files {
		fn := filepath.Join(p, f.Name())
		// Dotfile, like .git or whatever. Skip.
		if strings.HasPrefix(filepath.Base(fn), ".") {
			continue
		}

		if f.IsDir() {
			subdirs = append(subdirs, fn)
		} else if strings.HasSuffix(fn, ".yaml") {
			items = append(items, fn)
		}
	}

	cs := make(chan *Repo, len(subdirs)) // Sub-repo channel
	ci := make(chan Item, len(items))    // item channel

	// Start parsing subrepos
	for _, dir := range subdirs {
		go func(cs chan<- *Repo, dir string) {
			nr := NewRepo(dir)
			cs <- nr
		}(cs, dir)
	}

	// Start parsing items
	for _, fn := range items {
		go func(ci chan<- Item, fn string) {
			ni := r.loadItem(fn)
			ci <- ni
		}(ci, fn)
	}

	// Drain the items first
	for x := 0; x < len(items); x++ {
		item := <-ci
		// Control files start with an underscore and should not be stored as
		// normal Item documents.
		path := item.Path()
		id := item.ID()
		if strings.HasPrefix(asKey(path), "_") {
			r.Control[id] = item
		} else {
			r.Items[id] = item
		}
	}

	// And then drain the subrepos
	for x := 0; x < len(subdirs); x++ {
		sub := <-cs
		r.Subrepos[sub.Key] = sub
	}

	return &r
}

// ListRepos prints a sorted list of available repostiories.
func ListRepos(repos map[string]Repo) {
	keys := make([]string, 0, len(repos))
	for key := range repos {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key)
	}
}

// Keys returns a sorted list of the info keys in the repository
func (r *Repo) Keys() []string {
	keys := make([]string, 0, len(r.Items))
	for _, item := range r.Items {
		keys = append(keys, item.ID())
	}

	sort.Strings(keys)

	return keys
}

// SubrepoKeys returns a sorted list of the subrepo keys in the repository
func (r *Repo) SubrepoKeys() []string {
	keys := make([]string, 0, len(r.Subrepos))
	for _, sub := range r.Subrepos {
		keys = append(keys, sub.Key)
	}

	sort.Strings(keys)

	return keys
}

// GetHost will return a Host as defined by the list of arguments
//
// `args` is to be a string containing space separated identifiers to find a
// host category.
func (r *Repo) GetHost(def string) (h *Host) {
	args := strings.Split(def, " ")
	if len(args) < 2 {
		log.Fatal("Too few identifiers in host string. Need at least 2.")
	}

	args = append([]string{"hosts"}, args...)

	item, remaining, err := r.GetItem(args)
	host := item.(HostInfo)

	if err != nil {
		log.Print(err)
		log.Fatal("No host could be found")
	}

	cat := host.Types[remaining[0]]
	return cat.PrimaryHost()
}

// GetItem will return an Info as defined by the list of arguments
//
// If successful, a *Info is returned along with the remaining unparsed arguments.
func (r *Repo) GetItem(args []string) (Item, []string, error) {
	var item Item
	var ok bool

	repo, remaining, err := r.GetSubrepo(args)
	if err != nil {
		log.Print(repo)
		return nil, []string{}, errors.New(
			"No matching Info found because no subrepo matched the query.",
		)
	}

	if item, ok = repo.Items[remaining[0]]; ok {
		return item, remaining[1:], nil
	}

	return nil, []string{}, errors.New("No matching Info found.")
}

// GetSubrepo will return an Info as defined by the list of arguments
//
// If successful, a *Repo is returned along with the remaining unparsed arguments.
func (r *Repo) GetSubrepo(args []string) (*Repo, []string, error) {
	var err error

	if len(args) == 0 {
		return r, args, nil
	}

	arg := args[0]
	if repo, ok := r.Subrepos[arg]; ok {
		return repo.GetSubrepo(args[1:])
	}

	if _, ok := r.Items[arg]; !ok {
		err = fmt.Errorf("Subrepo did not exist: %s", arg)
	}
	return r, args, err
}

// ParentRepo parses the repo tree upwards until it finds the root repository
//
// This is used by things like command execution, where the current repository would be
// `commands` or a subrepository, but the root is needed for host discovery.
func (r *Repo) ParentRepo() *Repo {
	if &r.Parent == nil {
		return r
	}
	return r.Parent
}

// MakeCLI generates a cli.Command chain based on the repository structure
func (r *Repo) MakeCLI() (c cli.Command) {
	c = cli.Command{
		Name:     r.Key,
		Usage:    r.Summary,
		HideHelp: true,
	}

	// Make a list of subcommands to add into the Command.
	subcommands := make([]cli.Command, 0, len(r.Items)+len(r.Subrepos))

	// Loop over the subrepositories first, making sure that they are on top.
	for _, key := range r.SubrepoKeys() {
		subrepo := r.Subrepos[key]
		subcommands = append(subcommands, subrepo.MakeCLI())
	}

	// Then loop the item files.
	for _, key := range r.Keys() {
		item := r.Items[key]

		sc := cli.Command{
			Name:     item.ID(),
			Usage:    item.Summary(),
			HideHelp: true,
			Action:   item.Execute,
		}

		sc.Subcommands = append(sc.Subcommands, item.MakeCLI()...)

		subcommands = append(subcommands, sc)
	}

	c.Subcommands = subcommands

	return
}

func (r *Repo) loadItem(path string) Item {
	info, err := LoadItem(r, path)
	if err != nil {
		log.Println("Failed to load info: ", err)
	}
	return info
}

// Helper to run git commands inside of a repository
func (r *Repo) git(args ...string) {
	git(r.root, args...)
}
