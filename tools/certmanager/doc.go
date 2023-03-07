package main

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/pkg/doc"
	"os"
)

func main() {
	dir := os.Getenv("DOCS_OUT_CERT_MANAGER")
	if dir == "" {
		fmt.Printf("DOCS_OUT_CERT_MANAGER environment variable not set.\n")
		os.Exit(1)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error getting directory: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands.GetRootCmd().SubCommands() {
		// Find certificate-manager command
		if cmd.Command != nil && cmd.Command.Parent() != nil && cmd.Name() == "certificate-manager" {
			for _, subCmd := range cmd.SubCommands() {
				err := doc.WriteDocs(subCmd, dir)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}
}
