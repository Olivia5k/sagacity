package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)

func testHostInfo() *HostInfo {
	p := "test/repos/host_tests/printout/hosts/db.yaml"
	h := HostInfo{id: asKey(p), path: p, repo: &Repo{}}

	data, _ := ioutil.ReadFile(p)
	yaml.Unmarshal(data, &h)
	return &h
}

func TestHostInfoTypesAreSet(t *testing.T) {
	assert := assert.New(t)
	h := testHostInfo()

	assert.Equal(5, len(h.Types))
}

func TestHostInfoHostsAreSet(t *testing.T) {
	assert := assert.New(t)
	h := testHostInfo()

	assert.Equal(1, len(h.Types["master"].Hosts))
	assert.Equal(4, len(h.Types["ro"].Hosts))
	assert.Equal(1, len(h.Types["wal"].Hosts))
	assert.Equal(2, len(h.Types["standby"].Hosts))
	assert.Equal(2, len(h.Types["task"].Hosts))
}

func ExampleHostType() {
	data, err := ioutil.ReadFile("test/host_example_config.yaml")
	if err != nil {
		log.Fatal("No host configuration file found: ", err)
	}

	conf := &Config{}
	yaml.Unmarshal(data, conf)
	repos := LoadRepos(conf)

	app := BuildCLI(repos, conf)
	app.Run([]string{"sagacity", "printout", "hosts", "db"})

	// Output: [36;1mmaster[0m:
	//   Master database, read/write
	//   [33m[[0m[93;1m0[0m[33m][0m [34;1mdb1.cluster6.company.net[0m ([32;1mprimary[0m)

	// [36;1mro[0m:
	//   Read-only slaves
	//   [33m[[0m[93;1m0[0m[33m][0m [34;1mdb2.cluster3.company.net[0m
	//   [33m[[0m[93;1m1[0m[33m][0m [34;1mdb5.cluster3.company.net[0m
	//   [33m[[0m[93;1m2[0m[33m][0m [34;1mdb6.cluster3.company.net[0m
	//   [33m[[0m[93;1m3[0m[33m][0m [34;1mdb4.cluster3.company.net[0m ([32;1mprimary[0m) ([37mDesignated for long queries[0m)

	// [36;1mstandby[0m:
	//   Hot standby machines
	//   [33m[[0m[93;1m0[0m[33m][0m [34;1mdb8.cluster3.company.net[0m ([32;1mprimary[0m)
	//   [33m[[0m[93;1m1[0m[33m][0m [34;1mdb1.cluster3.company.net[0m ([37mHot standby, disaster recovery only[0m)

	// [36;1mtask[0m:
	//   task-only db machines
	//   [33m[[0m[93;1m0[0m[33m][0m [34;1mtaskdb1.cluster6.company.net[0m
	//   [33m[[0m[93;1m1[0m[33m][0m [34;1mtaskdb2.cluster6.company.net[0m

	// [36;1mwal[0m:
	//   WAL archive storage machines
	//   [33m[[0m[93;1m0[0m[33m][0m [34;1mdb7.cluster3.company.net[0m
}
