package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/pkg/doc"
)

const (
	DocsOutFolder = "docs/subcommands"
)

func main() {
	if err := os.MkdirAll(DocsOutFolder, os.ModePerm); err != nil {
		die.Die(fmt.Errorf("error creating directories: %w", err).Error())
	}

	err := doc.WriteDocs(commands.GetRootCmd(), DocsOutFolder)
	if err != nil {
		die.Die(fmt.Errorf("error writing command docs %w", err).Error())
	}

	err = doc.GenerateSummary(filepath.Join(DocsOutFolder, ".."))
	if err != nil {
		die.Die(fmt.Errorf("error writing summary %w", err).Error())
	}
}
