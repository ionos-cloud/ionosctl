package main

import (
	"github.com/ionos-cloud/ionosctl/commands"
)

var version = "master"

func main() {
	commands.Execute(version)
}
