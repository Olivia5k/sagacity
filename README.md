# sagacity

`sagacity` is a command line based knowledge base. Given `yaml` files it can
help you:

* Show and execute commands and tips.
* Open `ssh` connections to servers based on Puppet roles.
* Create and run dependency-graph based runsheets.

## Installation

`go get -u github.com/thiderman/sagacity`

## Usage

* `saga repo <install|update>`
Manage the repositories containing `yaml` recipes.

* `saga search <phrase> [<phrase>, ...]`
Search through the knowledge base and print matching

* `saga <til|random>`
Show a random article. Put this in your `$SHELL.rc`!

The rest of the usage is depending on what repositories you have
installed. The repositories contain definitions of commands

## License
MIT. See the LICENSE file.
