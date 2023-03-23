package main

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/pkg/doc"
)

const (
	DocsOutFolder = "docs"
)

func main() {
	if _, err := os.Stat(DocsOutFolder); err != nil {
		die.Die(fmt.Errorf("error getting directory stat: %w", err).Error())
	}

	err := doc.WriteDocs(commands.GetRootCmd(), DocsOutFolder)
	if err != nil {
		die.Die(fmt.Errorf("error writing command docs %w", err).Error())
	}

	err = doc.GenerateSummary(DocsOutFolder)
	if err != nil {
		die.Die(fmt.Errorf("error writing summary %w", err).Error())
	}
}
