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
	"sort"
	"strings"
	"sync"
)

// Repo represents a repository of information yaml files.
type Repo struct {
	Key      string `yaml:"key"`
	Summary  string `yaml:"summary"`
	Alias    string `yaml:"alias"`
	Info     map[string]Info
	Control  map[string]Info
	Subrepos map[string]Repo
	Parent   *Repo
	root     string
	wg       sync.WaitGroup
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

		if _, err := os.Stat(filepath.Join(fn, "_repo.yaml")); os.IsNotExist(err) {
			log.Println(fmt.Sprintf("Skipping repo %s: no _repo.yaml found.", file.Name()))
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			repos[asKey(fn)] = NewRepo(fn)
		}()
	}

	wg.Wait()
	return
}

// UpdateRepos will run git pull on the repos
func UpdateRepos(repos map[string]Repo) {
	for key, repo := range repos {
		log.Printf("Updating %s...", key)
		repo.git("pull", "origin", "master")
	}
}

// AddRepo clones a new repository
func AddRepo(root, name, url string) {
	dir := filepath.Join(root, name)
	git("", "clone", url, dir)
	log.Print("Repository added!")
}

// NewRepo loads a repository on a path
func NewRepo(p string) (r Repo) {
	p = getPath(p)
	r = Repo{Key: asKey(p), root: p}

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

	r.Info = make(map[string]Info)
	r.Control = make(map[string]Info)
	r.Subrepos = make(map[string]Repo)

	filepath.Walk(r.root, r.walk)
	r.wg.Wait()

	return
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
	keys := make([]string, 0, len(r.Info))
	for _, info := range r.Info {
		keys = append(keys, info.ID)
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

	info, remaining, err := r.GetInfo(args)
	if err != nil {
		log.Print(err)
		log.Fatal("No host could be found")
	}

	cat := info.Hosts[remaining[0]]
	return cat.PrimaryHost()
}

// GetInfo will return an Info as defined by the list of arguments
//
// If successful, a *Info is returned along with the remaining unparsed arguments.
func (r *Repo) GetInfo(args []string) (*Info, []string, error) {
	var info Info
	var ok bool

	repo, remaining, err := r.GetSubrepo(args)
	if err != nil {
		log.Print(repo)
		return nil, []string{}, errors.New(
			"No matching Info found because no subrepo matched the query.",
		)
	}

	if info, ok = repo.Info[remaining[0]]; ok {
		return &info, remaining[1:], nil
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

	if _, ok := r.Info[arg]; !ok {
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
	subcommands := make([]cli.Command, 0, len(r.Info)+len(r.Subrepos))

	// Loop over the subrepositories first, making sure that they are on top.
	for _, key := range r.SubrepoKeys() {
		subrepo := r.Subrepos[key]
		subcommands = append(subcommands, subrepo.MakeCLI())
	}

	// Then loop the info files.
	for _, key := range r.Keys() {
		info := r.Info[key]

		sc := cli.Command{
			Name:     info.ID,
			Usage:    info.Summary,
			HideHelp: true,
			Action: func(c *cli.Context) {
				info.Execute(r, c)
			},
		}

		if info.Type == "host" {
			sc.Subcommands = append(sc.Subcommands, MakeHostCLI(&info)...)
		} else if info.Type == "command" {
			sc.Subcommands = append(sc.Subcommands, MakeCommandCLI(&info)...)
		}

		subcommands = append(subcommands, sc)
	}

	c.Subcommands = subcommands

	return
}

func (r *Repo) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("walk error: ", err)
		return err
	}

	// Dotfile, like .git or whatever. Skip.
	if strings.HasPrefix(filepath.Base(path), ".") {
		return filepath.SkipDir
	}

	if info.IsDir() && r.isSubrepo(path) {
		r.wg.Add(1)
		go r.loadSubrepo(path)

		// Return SkipDir since the directory will be parsed by the
		// NewRepo call inside of loadSubrepo()
		return filepath.SkipDir

	} else if strings.HasSuffix(path, ".yaml") {
		r.wg.Add(1)
		go r.loadInfo(path)
	}

	return nil
}

func (r *Repo) loadInfo(path string) {
	defer r.wg.Done()

	info, err := LoadInfo(r, path)
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

func (r *Repo) loadSubrepo(path string) {
	defer r.wg.Done()
	nr := NewRepo(path)
	nr.Parent = r
	r.Subrepos[nr.Key] = nr
}

func (r *Repo) isSubrepo(path string) bool {
	// This is the root...
	if r.root == path {
		return false
	}

	return true
}

// Helper to run git commands inside of a repository
func (r *Repo) git(args ...string) {
	git(r.root, args...)
}
