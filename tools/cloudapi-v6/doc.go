package main

import (
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/commands"
	"github.com/ionos-cloud/ionosctl/tools/internal"
)

func main() {
	dir := os.Getenv("DOCS_OUT_V6")
	if dir == "" {
		fmt.Printf("DOCS_OUT_V6 environment variable not set.\n")
		os.Exit(1)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error getting directory: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands.GetRootCmd().SubCommands() {
		if cmd.Name() != "dbaas-pgsql" && cmd.Name() != "version" && cmd.Name() != "login" {
			err := internal.WriteDocs(cmd, dir)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
