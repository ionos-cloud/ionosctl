package main

import (
	"fmt"
<<<<<<< HEAD
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
=======
	"github.com/ionos-cloud/ionosctl/commands"
	"github.com/ionos-cloud/ionosctl/pkg/doc"
	"os"
>>>>>>> ca183e9b (doc: doc update script with refactor for less duplicated code)
)

func main() {
	dir := os.Getenv("DOCS_OUT_DBAAS_POSTGRES")
	if dir == "" {
		fmt.Printf("DOCS_OUT_DBAAS_POSTGRES environment variable not set.\n")
		os.Exit(1)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error getting directory: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands.GetRootCmd().SubCommands() {
		// Find dbaas command
		if cmd.Command != nil && cmd.Command.Parent() != nil && cmd.Name() == "dbaas" {
			for _, subCmd := range cmd.SubCommands() {
				// Find postgres command
				if subCmd.Name() == "postgres" {
					err := doc.WriteDocs(subCmd, dir)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						os.Exit(1)
					}
				}
			}
		}
	}
}
