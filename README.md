# sagacity

`sagacity` is a command line based knowledge base and execution helper.

### Current features
* Open `ssh` connections to servers based on their roles.

### Coming features
* Show and execute commands on single or multiple hosts.
* Ping hosts for heart metrics.
* Create and run dependency-graph based runsheets.

## Installation

`go get -u github.com/thiderman/sagacity`

For good UX in `bash` and `zsh` add the following to your shell rc:

```
PROG=sagacity source $GOPATH/src/github.com/codegangsta/cli/autocomplete/$(basename $SHELL)_autocomplete
alias sp=sagacity
```

## Usage

* `sagacity repo <add|update>`
Manage the repositories containing `yaml` recipes.

## License
MIT. See the LICENSE file.
