package main

import (
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/commands"
	"github.com/ionos-cloud/ionosctl/tools/internal"
)

const rootDbaasPgsqlCmdName = "dbaas-pgsql"

func main() {
	dir := os.Getenv("DOCS_OUT_DBAAS_PGSQL")
	if dir == "" {
		fmt.Printf("DOCS_OUT_DBAAS_PGSQL environment variable not set.\n")
		os.Exit(1)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error getting directory: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands.GetRootCmd().SubCommands() {
		if cmd.Command.Name() == rootDbaasPgsqlCmdName {
			err := internal.WriteDocs(cmd, dir)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
